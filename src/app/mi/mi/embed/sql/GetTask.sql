SELECT
    TaskID,
    CreatedTime
FROM
    Task
WHERE
    TaskID = '%s'
LIMIT
    1;