SELECT
    LimitID,
    TaskID,
    UpdatedTime,
    LimitTime
FROM
    LimitInfo
WHERE
    TaskID = '%s'
ORDER BY
    UpdatedTime DESC
LIMIT
    1;