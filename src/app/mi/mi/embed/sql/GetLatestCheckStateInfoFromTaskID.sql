SELECT
    CheckStateID,
    TaskID,
    UpdatedTime,
    IsChecked
FROM
    CheckStateInfo
WHERE
    TaskID = '%s'
ORDER BY
    UpdatedTime DESC
LIMIT
    1;