SELECT
    BoardInfoID,
    TaskID,
    UpdatedTime,
    BoardName
FROM
    BoardInfo
WHERE
    TaskID = '%s'
ORDER BY
    UpdatedTime DESC
LIMIT
    1;