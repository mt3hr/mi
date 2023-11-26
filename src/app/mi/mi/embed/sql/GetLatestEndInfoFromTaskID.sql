SELECT
    EndID,
    TaskID,
    UpdatedTime,
    EndTime
FROM
    EndInfo
WHERE
    TaskID = '%s'
ORDER BY
    UpdatedTime DESC
LIMIT
    1;