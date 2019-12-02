import sys, getopt
import torch
import time
import numpy as np
from sklearn.model_selection import train_test_split
from torch.autograd import Variable
from torch.nn import functional as F
import cx_Oracle

def writeOutput(key, value):
    separator = ":"
    print(key + separator + str(value))
    sys.stdout.flush()

def selectStationStatisticsByDay(db):
    cur = db.cursor()
    statement = """
        SELECT  
            citibike_station.latitude,
            citibike_station.longitude,
            citibike_station_statistics_by_day.number_arrivals_previous_90_days,
            citibike_station_statistics_by_day.number_departures_previous_90_days,
            TO_CHAR(citibike_station_statistics_by_day.m_date, 'DAY', 'NLS_DATE_LANGUAGE=''numeric date language'''),
            TO_CHAR(citibike_station_statistics_by_day.m_date, 'MONTH', 'NLS_DATE_LANGUAGE=''numeric date language'''),
            citibike_station_statistics_by_day.number_arrivals
        FROM citibike_station
        JOIN citibike_station_statistics_by_day ON citibike_station.station_id = citibike_station_statistics_by_day.station_id
        WHERE citibike_station_statistics_by_day.m_date >= to_date('29-10-2013', 'dd-mm-yyyy')
        AND citibike_station_statistics_by_day.number_arrivals_previous_90_days > 0
        AND citibike_station_statistics_by_day.number_departures_previous_90_days > 0
    """
    cur.execute(statement)
    res = cur.fetchall()
    if len(res) == 0:
        return []

    npRes = np.array(res).astype(np.float32)
    #this line clears out nan values, need to fix in DB
    npRes =  npRes[~np.isnan(npRes).any(axis=1)]
    npRes =  npRes[~(npRes==0).any(axis=1)]
    x_data = npRes[:, :6].astype(np.float32)
    y_data = npRes[:,6].astype(np.float32)
    return x_data, y_data

def main(argv):
    device = torch.device("cuda:0")
    opts, args = getopt.getopt(argv, 'e:')

    epochs = 50000
    for opt, arg in opts:
        if opt == '-e':
            epochs = int(arg)

    writeOutput("scriptName", "Citibike Station Arrivals By Day")
    writeOutput("processor", "GPU")
    writeOutput("step", "Connecting to Oracle DB")
    db = cx_Oracle.connect(user="ADMIN", password="Oracle12345!", dsn="burlmigration_high")
    print("Connected to Oracle ADW")

    startTime = time.time()
    writeOutput("startTime", startTime)

    writeOutput("step", "Querying data")
    sqlStartTime = time.time()
    writeOutput("sqlStartTime", sqlStartTime)
    x_data, y_data = selectStationStatisticsByDay(db)
    writeOutput("sqlTime", "{:.3f}".format(time.time() - sqlStartTime))

    writeOutput("step", "Preparing model")
    pyTorchModelStartTime = time.time()
    writeOutput("pyTorchModelStartTime", pyTorchModelStartTime)

    x_norm = x_data / x_data.max(axis=0)
    y_norm = y_data / y_data.max(axis=0)

    x_train, x_test, y_train, y_test = train_test_split(x_norm, y_norm, test_size=0.20, random_state=42)


    input_dim = len(x_data[0])
    output_dim = 1
    writeOutput("dataLength", len(x_data))
    writeOutput("dataWidth", input_dim)

    print("Loading data and model on to GPU")

    class LogisticRegression(torch.nn.Module):

        def __init__(self, input_dim, output_dim):

            super(LogisticRegression, self).__init__() 
            # Calling Super Class's constructor
            self.linear = torch.nn.Linear(input_dim, output_dim)
            # nn.linear is defined in nn.Module

        def forward(self, x):
            # Here the forward pass is simply a linear function

            out = torch.sigmoid(self.linear(x))
            return out



    x_tensor = torch.Tensor(x_norm)
    y_tensor = torch.Tensor(y_norm)
    y_ok = y_tensor.unsqueeze(1)

    trainds = torch.utils.data.TensorDataset(x_tensor, y_ok)
    trainloader = torch.utils.data.DataLoader(trainds, batch_size=128, shuffle=False, num_workers=1000)

    model = LogisticRegression(input_dim,output_dim).to(device)
    criterion = torch.nn.MSELoss().to(device)# Mean Squared Loss
    l_rate = 0.5
    optimizer = torch.optim.SGD(model.parameters(), lr = l_rate) #Stochastic Gradient Descent

    writeOutput("pyTorchModelTime", "{:.3f}".format(time.time() - pyTorchModelStartTime))

    writeOutput("step", "Training Model")
    writeOutput("epochs", epochs)
    trainingStartTime = time.time()
    writeOutput("trainingStartTime", trainingStartTime)

    model.train()
    for epoch in range(epochs):
        print(epoch)
        running_loss = 0.0
        for i, data in enumerate(trainloader, 0):
            inputs, labels = data[0].to(device), data[1].to(device)
            optimizer.zero_grad()
            y_pred = model(inputs)
            loss = criterion(y_pred, labels)
            loss.backward()
            optimizer.step()
            running_loss += loss.item()
    #for epoch in range(epochs):
    #    y_pred = model(x_tensor)
    #    loss = criterion(y_pred, y_ok)
    #    optimizer.zero_grad()
    #    loss.backward()
    #    optimizer.step()
        if (epoch % 100 == 0):
            pctComplete = epoch / epochs * 100
            writeOutput("percentComplete", "{:.2f}".format(pctComplete))
            writeOutput("loss", "{:.6f}".format(running_loss))
            print(running_loss)
            
    writeOutput("totalTime", "{:.3f}".format(time.time() - startTime))
    writeOutput("trainingTime", "{:.3f}".format(time.time() - trainingStartTime))

    writeOutput("step", "Test")
    model.eval()
    diffs = []
    for i, val in enumerate(x_test):
        x_test_tensor = torch.from_numpy(val).to(device)
        predicted = model(x_test_tensor)
        diff = abs(predicted.item() - y_test[i])
        diffs.append(float(diff))
    avg = np.average(diffs)
    writeOutput("accuracy", avg * 100)
    writeOutput("step", "Finished")
    print('fin')

if __name__ == "__main__":
   main(sys.argv[1:])