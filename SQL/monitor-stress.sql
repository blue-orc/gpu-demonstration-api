CREATE TABLE peak_iops_measurement (capture_timestamp date, total_read_io number, total_write_io number, total_io number, total_read_bytes number, total_write_bytes number, total_bytes number);


DECLARE

run_duration number := 120;
capture_gap number := 5;
loop_count number :=run_duration/capture_gap;

rdio number; 
wtio number;
prev_rdio number :=0; 
prev_wtio number :=0;

rdbt number; 
wtbt number;
prev_rdbt number; 
prev_wtbt number;

BEGIN
    FOR i in 1..loop_count LOOP
    
    SELECT SUM(value) INTO rdio from gv$sysstat 
    WHERE name ='physical read total IO requests';
    
    SELECT SUM(value) INTO wtio from gv$sysstat 
    WHERE name ='physical write total IO requests';
    
    SELECT SUM(value)* 0.000008 INTO rdbt from gv$sysstat 
    WHERE name ='physical read total bytes'; 
    
    SELECT SUM(value* 0.000008) INTO wtbt from gv$sysstat 
    WHERE name ='physical write total bytes';
    
    IF i > 1  THEN
        INSERT INTO peak_iops_measurement (capture_timestamp, total_read_io, total_write_io, total_io, total_read_bytes, total_write_bytes, total_bytes)
        VALUES (sysdate,(rdio-prev_rdio)/5,(wtio-prev_wtio)/5,((rdio-prev_rdio)/5)+((wtio-prev_wtio))/5,(rdbt-prev_rdbt)/5,(wtbt-prev_wtbt)/5,((rdbt-prev_rdbt)/5)+((wtbt-prev_wtbt))/5);
    END IF;
    
    prev_rdio := rdio;
    prev_wtio := wtio;
    
    prev_rdbt := rdbt;
    prev_wtbt := wtbt;
    
    DBMS_LOCK.SLEEP(capture_gap);
    
    END LOOP;
    
    COMMIT;
    EXCEPTION  
        WHEN OTHERS THEN ROLLBACK;
END;