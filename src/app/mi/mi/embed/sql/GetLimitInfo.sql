SELECT
    LimitID,
    TaskID,
    UpdatedTime,
    LimitTime
FROM
    LimitInfo
WHERE
    LimitID = '%s'
LIMIT
    1;