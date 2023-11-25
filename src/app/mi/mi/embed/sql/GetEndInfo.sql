SELECT
    EndID,
    TaskID,
    UpdatedTime,
    EndTime
FROM
    EndInfo
WHERE
    EndID = '%s'
LIMIT
    1;