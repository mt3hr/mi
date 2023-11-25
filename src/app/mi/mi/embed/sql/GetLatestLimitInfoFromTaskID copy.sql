SELECT
    StartID,
    TaskID,
    UpdatedTime,
    LimitTime
FROM
    StartInfo
WHERE
    TaskID = '%s'
ORDER BY
    UpdatedTime DESC
LIMIT
    1;