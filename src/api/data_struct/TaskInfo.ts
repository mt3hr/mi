import BoardInfo from "./BoardInfo";
import CheckStateInfo from "./CheckStateInfo";
import LimitInfo from "./LimitInfo";
import Task from "./Task";
import TaskTitleInfo from "./TaskTitleInfo";

export default class TaskInfo {
    public task: Task = new Task()
    public task_title_info: TaskTitleInfo = new TaskTitleInfo()
    public check_state_info: CheckStateInfo = new CheckStateInfo()
    public limit_info: LimitInfo = new LimitInfo()
    public board_info: BoardInfo = new BoardInfo()
}