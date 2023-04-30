import AddTagRequest from "./AddTagRequest";
import AddTagResponse from "./AddTagResponse";
import DeleteTaskRequest from "./DeleteTaskRequest";
import DeleteTaskResponse from "./DeleteTaskResponse";
import GetBoardStructRequest from "./GetBoardStructRequest";
import GetBoardStructResponse from "./GetBoardStructResponse";
import GetTagStructRequest from "./GetTagStructRequest";
import GetTagStructResponse from "./GetTagStructResponse";
import GetTasksFromBoardRequest from "./GetTasksFromBoardRequest";
import GetTasksFromBoardResponse from "./GetTasksFromBoardResponse";
import GetTaskRequest from "./GetTaskRequest";
import GetTaskResponse from "./GetTaskResponse";
import UpdateTaskRequest from "./UpdateTaskRequest";
import UpdateTaskResponse from "./UpdateTaskResponse";
import AddTextRequest from "./AddTextRequest";
import AddTextResponse from "./AddTextResponse";
import GetTagsRelatedTaskRequest from "./GetTagsRelatedTaskRequest";
import GetTagsRelatedTaskResponse from "./GetTagsRelatedTaskResponse";
import GetTextsRelatedTaskRequest from "./GetTextsRelatedTaskRequest";
import GetTextsRelatedTaskResponse from "./GetTextsRelatedTaskResponse";
import DeleteTagRequest from "./DeleteTagRequest";
import DeleteTagResponse from "./DeleteTagResponse";
import DeleteTextRequest from "./DeleteTextRequest";
import DeleteTextResponse from "./DeleteTextResponse";

export default class MiServerAPI {
    public get_board_struct(get_board_struct_request: GetBoardStructRequest): GetBoardStructResponse {
        return new GetBoardStructResponse()//TODO
    }
    public get_tag_struct(get_tag_struct_request: GetTagStructRequest): GetTagStructResponse {
        return new GetTagStructResponse()//TODO
    }
    public add_task(add_task_request: AddTagRequest): AddTagResponse {
        return new AddTagResponse()//TODO
    }
    public update_task(update_task_request: UpdateTaskRequest): UpdateTaskResponse {
        return new UpdateTaskResponse()//TODO
    }
    public delete_task(delete_task_request: DeleteTaskRequest): DeleteTaskResponse {
        return new DeleteTaskResponse()//TODO
    }
    public get_task(get_task_request: GetTaskRequest): GetTaskResponse {
        return new GetTaskResponse()//TODO
    }
    public get_tasks_from_board(get_tasks_from_board_request: GetTasksFromBoardRequest): GetTasksFromBoardResponse {
        return new GetTasksFromBoardResponse()//TODO
    }
    public add_tag(add_tag_request: AddTagRequest): AddTagResponse {
        return new AddTagResponse();//TODO
    }
    public add_text(add_text_request: AddTextRequest): AddTextResponse {
        return new AddTextResponse()//TODO
    }
    public get_tags_related_task(get_tags_related_task_request: GetTagsRelatedTaskRequest): GetTagsRelatedTaskResponse {
        return new GetTagsRelatedTaskResponse()//TODO
    }
    public get_texts_related_task(get_texts_related_task_request: GetTextsRelatedTaskRequest): GetTextsRelatedTaskResponse {
        return new GetTextsRelatedTaskResponse()//TODO
    }
    public delete_tag(delete_tag_request: DeleteTagRequest): DeleteTagResponse {
        return new DeleteTagResponse()//TODO
    }
    public delete_text(delete_text_request: DeleteTextRequest): DeleteTextResponse {
        return new DeleteTextResponse()//TODO
    }
}