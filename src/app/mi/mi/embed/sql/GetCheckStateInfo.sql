SELECT
    CheckStateID,
    TaskID,
    UpdatedTime,
    IsChecked
FROM
    CheckStateInfo
WHERE
    CheckStateID = '%s'
LIMIT
    1;