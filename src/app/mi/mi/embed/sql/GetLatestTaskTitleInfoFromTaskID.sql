SELECT
    TaskTitleID,
    TaskID,
    UpdatedTime,
    Title
FROM
    TaskTitleInfo
WHERE
    TaskID = '%s'
ORDER BY
    UpdatedTime DESC
LIMIT
    1;