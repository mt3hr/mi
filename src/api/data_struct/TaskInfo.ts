import BoardInfo from "./BoardInfo";
import CheckStateInfo from "./CheckStateInfo";
import EndInfo from "./EndInfo";
import LimitInfo from "./LimitInfo";
import StartInfo from "./StartInfo";
import Task from "./Task";
import TaskTitleInfo from "./TaskTitleInfo";

export default class TaskInfo {
    public task: Task = new Task()
    public task_title_info: TaskTitleInfo = new TaskTitleInfo()
    public check_state_info: CheckStateInfo = new CheckStateInfo()
    public limit_info: LimitInfo = new LimitInfo()
    public start_info: StartInfo = new StartInfo()
    public end_info: EndInfo = new EndInfo()
    public board_info: BoardInfo = new BoardInfo()
}