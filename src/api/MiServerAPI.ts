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
import AddTaskResponse from "./AddTaskResponse";

const get_board_struct_address = "/api/board_struct"
const get_tag_struct_address = "/api/tag_struct"
const add_task_address = "/api/task"
const update_task_address = "/api/task"
const delete_task_address = "/api/task"
const get_task_address = "/api/task"
const get_tasks_from_board_address = "/api/board_task"
const add_tag_address = "/api/tag"
const add_text_address = "/api/text"
const get_tags_related_task_address = "/api/task_tag"
const get_texts_related_task_address = "/api/task_text"
const delete_tag_address = "/api/tag"
const delete_text_address = "/api/tag"
const get_board_struct_method = "get"
const get_tag_struct_method = "get"
const add_task_method = "post"
const update_task_method = "put"
const delete_task_method = "delete"
const get_task_method = "get"
const get_tasks_from_board_method = "get"
const add_tag_method = "post"
const add_text_method = "post"
const get_tags_related_task_method = "get"
const get_texts_related_task_method = "get"
const delete_tag_method = "delete"
const delete_text_method = "delete"

export default class MiServerAPI {
    public async get_board_struct(get_board_struct_request: GetBoardStructRequest): Promise<GetBoardStructResponse> {
        const res = await fetch(get_board_struct_address, {
            method: get_board_struct_method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(get_board_struct_request),
        })
        const json = await res.json()
        const response: GetBoardStructResponse = json
        return response
    }
    public async get_tag_struct(get_tag_struct_request: GetTagStructRequest): Promise<GetTagStructResponse> {
        const res = await fetch(get_tag_struct_address, {
            method: get_tag_struct_method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(get_tag_struct_request),
        })
        const json = await res.json()
        const response: GetTagStructResponse = json
        return response
    }
    public async add_task(add_task_request: AddTagRequest): Promise<AddTagResponse> {
        const res = await fetch(add_task_address, {
            method: add_task_method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(add_task_request),
        })
        const json = await res.json()
        const response: AddTaskResponse = json
        return response
    }
    public async update_task(update_task_request: UpdateTaskRequest): Promise<UpdateTaskResponse> {
        const res = await fetch(update_task_address, {
            method: update_task_method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(update_task_request),
        })
        const json = await res.json()
        const response: UpdateTaskResponse = json
        return response
    }
    public async delete_task(delete_task_request: DeleteTaskRequest): Promise<DeleteTaskResponse> {
        const res = await fetch(delete_task_address, {
            method: delete_task_method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(delete_task_request),
        })
        const json = await res.json()
        const response: DeleteTaskResponse = json
        return response
    }
    public async get_task(get_task_request: GetTaskRequest): Promise<GetTaskResponse> {
        const res = await fetch(get_task_address, {
            method: get_task_method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(get_task_request),
        })
        const json = await res.json()
        const response: GetTaskResponse = json
        return response
    }
    public async get_tasks_from_board(get_tasks_from_board_request: GetTasksFromBoardRequest): Promise<GetTasksFromBoardResponse> {
        const res = await fetch(get_tasks_from_board_address, {
            method: get_tasks_from_board_method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(get_tasks_from_board_request),
        })
        const json = await res.json()
        const response: GetTasksFromBoardResponse = json
        return response
    }
    public async add_tag(add_tag_request: AddTagRequest): Promise<AddTagResponse> {
        const res = await fetch(add_tag_address, {
            method: add_tag_method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(add_tag_request),
        })
        const json = await res.json()
        const response: AddTagResponse = json
        return response
    }
    public async add_text(add_text_request: AddTextRequest): Promise<AddTextResponse> {
        const res = await fetch(add_text_address, {
            method: add_text_method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(add_text_request),
        })
        const json = await res.json()
        const response: AddTextResponse = json
        return response
    }
    public async get_tags_related_task(get_tags_related_task_request: GetTagsRelatedTaskRequest): Promise<GetTagsRelatedTaskResponse> {
        const res = await fetch(get_tags_related_task_address, {
            method: get_tags_related_task_method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(get_tags_related_task_request),
        })
        const json = await res.json()
        const response: GetTagsRelatedTaskResponse = json
        return response
    }
    public async get_texts_related_task(get_texts_related_task_request: GetTextsRelatedTaskRequest): Promise<GetTextsRelatedTaskResponse> {
        const res = await fetch(get_texts_related_task_address, {
            method: get_texts_related_task_method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(get_texts_related_task_request),
        })
        const json = await res.json()
        const response: GetTextsRelatedTaskResponse = json
        return response
    }
    public async delete_tag(delete_tag_request: DeleteTagRequest): Promise<DeleteTagResponse> {
        const res = await fetch(delete_tag_address, {
            method: delete_tag_method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(delete_tag_request),
        })
        const json = await res.json()
        const response: DeleteTagResponse = json
        return response
    }
    public async delete_text(delete_text_request: DeleteTextRequest): Promise<DeleteTextResponse> {
        const res = await fetch(delete_text_address, {
            method: delete_text_method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(delete_text_request),
        })
        const json = await res.json()
        const response: DeleteTextResponse = json
        return response
    }
}