import CheckState from "./CheckState"
import SortType from "./SortType"

export default class TaskSearchQuery {
    public board: string = ""
    public tags: Array<string> = new Array<string>()
    public word: string = ""
    public check_state: CheckState = CheckState.NoCheckOnly
    public sort_type: SortType = SortType.CreatedTimeDesc
}