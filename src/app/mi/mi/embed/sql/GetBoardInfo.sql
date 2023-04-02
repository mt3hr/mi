SELECT
    BoardInfoID,
    TaskID,
    UpdatedTime,
    BoardName
FROM
    BoardInfo
WHERE
    BoardInfoID = '%s'
LIMIT
    1;