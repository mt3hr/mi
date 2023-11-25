SELECT
    StartID,
    TaskID,
    UpdatedTime,
    StartTime
FROM
    StartInfo
WHERE
    StartID = '%s'
LIMIT
    1;