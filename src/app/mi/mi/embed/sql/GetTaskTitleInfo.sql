SELECT
    TaskTitleID,
    TaskID,
    UpdatedTime,
    Title
FROM
    TaskTitleInfo
WHERE
    TaskTitleID = '%s'
LIMIT
    1;