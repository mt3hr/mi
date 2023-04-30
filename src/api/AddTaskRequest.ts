import MiRequest from "./MiRequest";
import TaskInfo from "./data_struct/TaskInfo";

export default class AddTaskRequest extends MiRequest {
    public task_info: TaskInfo = new TaskInfo()
}