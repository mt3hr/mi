SELECT
    StartID,
    TaskID,
    UpdatedTime,
    StartTime
FROM
    StartInfo
WHERE
    TaskID = '%s'
ORDER BY
    UpdatedTime DESC
LIMIT
    1;