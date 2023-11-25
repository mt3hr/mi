SELECT
    EndID,
    TaskID,
    UpdatedTime,
    EndTime
FROM
    LimitInfo
WHERE
    TaskID = '%s'
ORDER BY
    UpdatedTime DESC
LIMIT
    1;