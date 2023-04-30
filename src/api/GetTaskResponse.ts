import MiResponse from "./MiResponse";
import TaskInfo from "./data_struct/TaskInfo";

export default class GetTaskResponse extends MiResponse {
    public task_info: TaskInfo = new TaskInfo()
}