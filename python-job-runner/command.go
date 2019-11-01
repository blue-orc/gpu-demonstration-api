package PythonJobRunner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type PytorchJobStatus struct {
	Status          string
	PercentComplete string
	Loss            string
	ExecutionTime   string
	DataLength      string
	DataWidth       string
	Epochs          string
}

var Status PytorchJobStatus

func GetStatusJSON() ([]byte, error) {
	sBytes, err := json.Marshal(Status)
	if err != nil {
		return sBytes, err
	}
	return sBytes, nil
}

func Run() {
	if Status.Status == "Running" {
		return
	}

	Status.Status = "Running"
	cmd := exec.Command("python3", "/home/ubuntu/go/src/gpu-demonstration-api/python-job-runner/scripts/pytorch-training.py")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Run Python Script Error: " + err.Error())
		Status.Status = "Finished"
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("Run Python Script Error: " + err.Error())
		Status.Status = "Finished"
		return
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println("Run Python Script Error: " + err.Error())
		Status.Status = "Finished"
		return
	}

	go updateStatus(stdout)
	go updateStatus(stderr)
	cmd.Wait()
	Status.Status = "Finished"
}

func updateStatus(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		txt := scanner.Text()
		res := strings.Split(txt, ":")
		if res[0] == "dataLength" {
			Status.DataLength = res[1]
		} else if res[0] == "dataWidth" {
			Status.DataWidth = res[1]
		} else if res[0] == "epochs" {
			Status.Epochs = res[1]
		} else if res[0] == "percentComplete" {
			Status.PercentComplete = res[1]
		} else if res[0] == "loss" {
			Status.Loss = res[1]
		} else if res[0] == "executionTime" {
			Status.ExecutionTime = res[1]
		} else {
			fmt.Println(scanner.Text())
		}
	}
}
