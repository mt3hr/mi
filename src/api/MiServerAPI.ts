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
import GetTagRequest from "./GetTagRequest";
import GetTagResponse from "./GetTagResponse";
import GetTextRequest from "./GetTextRequest";
import GetTextResponse from "./GetTextResponse";
import AddTaskRequest from "./AddTaskRequest";
import GetApplicationConfigRequest from "./GetApplicationConfigRequest";
import GetApplicationConfigResponse from "./GetApplicationConfigResponse";
import GetBoardNamesRequest from "./GetBoardNamesRequest";
import GetBoardNamesResponse from "./GetBoardNamesResponse";
import GetTagNamesRequest from "./GetTagNamesRequest";
import GetTagNamesResponse from "./GetTagNamesResponse";

const get_board_struct_address = "/api/get_board_struct"
const get_tag_struct_address = "/api/get_tag_struct"
const add_task_address = "/api/add_task"
const update_task_address = "/api/update_task"
const delete_task_address = "/api/delete_task"
const get_task_address = "/api/get_task"
const get_tasks_from_board_address = "/api/get_board_task"
const add_tag_address = "/api/add_tag"
const add_text_address = "/api/add_text"
const get_tags_related_task_address = "/api/get_task_tag"
const get_texts_related_task_address = "/api/get_task_text"
const get_tag_address = "/api/get_tag"
const get_text_address = "/api/get_text"
const delete_tag_address = "/api/delete_tag"
const delete_text_address = "/api/delete_text"
const get_tag_names_address = "/api/get_tag_names"
const get_board_names_address = "/api/get_board_names"
const get_application_config_address = "/api/get_application_config"
const get_board_struct_method = "post"
const get_tag_struct_method = "post"
const add_task_method = "post"
const update_task_method = "post"
const delete_task_method = "post"
const get_task_method = "post"
const get_tasks_from_board_method = "post"
const add_tag_method = "post"
const add_text_method = "post"
const get_tags_related_task_method = "post"
const get_texts_related_task_method = "post"
const delete_tag_method = "post"
const delete_text_method = "post"
const get_tag_method = "post"
const get_text_method = "post"
const get_tag_names_method = "post"
const get_board_names_method = "post"
const get_application_config_method = "post"

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
    public async add_task(add_task_request: AddTaskRequest): Promise<AddTagResponse> {
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
        response.task_info.board_info.updated_time = new Date(response.task_info.board_info.updated_time)
        response.task_info.check_state_info.updated_time = new Date(response.task_info.check_state_info.updated_time)
        response.task_info.limit_info.updated_time = new Date(response.task_info.limit_info.updated_time)
        response.task_info.limit_info.limit = response.task_info.limit_info.limit ? new Date(response.task_info.limit_info.limit) : null
        response.task_info.start_info.updated_time = new Date(response.task_info.start_info.updated_time)
        response.task_info.start_info.start = response.task_info.start_info.start ? new Date(response.task_info.start_info.start) : null
        response.task_info.end_info.updated_time = new Date(response.task_info.end_info.updated_time)
        response.task_info.end_info.end = response.task_info.end_info.end ? new Date(response.task_info.end_info.end) : null
        response.task_info.task.created_time = new Date(response.task_info.task.created_time)
        response.task_info.task_title_info.updated_time = new Date(response.task_info.task_title_info.updated_time)
        return response
    }
    public async get_tasks_from_board(get_tasks_from_board_request: GetTasksFromBoardRequest, abort_controller: AbortController): Promise<GetTasksFromBoardResponse> {
        const res = await fetch(get_tasks_from_board_address, {
            method: get_tasks_from_board_method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(get_tasks_from_board_request),
            signal: abort_controller.signal,
        })
        const json = await res.json()
        const response: GetTasksFromBoardResponse = json
        if (response.boards_tasks) {
            for (let i = 0; i < response.boards_tasks.length; i++) {
                response.boards_tasks[i].board_info.updated_time = new Date(response.boards_tasks[i].board_info.updated_time)
                response.boards_tasks[i].check_state_info.updated_time = new Date(response.boards_tasks[i].check_state_info.updated_time)
                response.boards_tasks[i].limit_info.updated_time = new Date(response.boards_tasks[i].limit_info.updated_time)
                response.boards_tasks[i].limit_info.limit = response.boards_tasks[i].limit_info.limit ? new Date(response.boards_tasks[i].limit_info.limit!) : null
                response.boards_tasks[i].start_info.updated_time = new Date(response.boards_tasks[i].start_info.updated_time)
                response.boards_tasks[i].start_info.start = response.boards_tasks[i].start_info.start ? new Date(response.boards_tasks[i].start_info.start!) : null
                response.boards_tasks[i].end_info.updated_time = new Date(response.boards_tasks[i].end_info.updated_time)
                response.boards_tasks[i].end_info.end = response.boards_tasks[i].end_info.end ? new Date(response.boards_tasks[i].end_info.end!) : null
                response.boards_tasks[i].task.created_time = new Date(response.boards_tasks[i].task.created_time)
                response.boards_tasks[i].task_title_info.updated_time = new Date(response.boards_tasks[i].task_title_info.updated_time)
            }
        }
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
        for (let i = 0; i < response.tags.length; i++) {
            response.tags[i].time = new Date(response.tags[i].time)
        }
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
        for (let i = 0; i < response.texts.length; i++) {
            response.texts[i].time = new Date(response.texts[i].time)
        }
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
    public async get_tag(get_tag_request: GetTagRequest): Promise<GetTagResponse> {
        const res = await fetch(get_tag_address, {
            method: get_tag_method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(get_tag_request),
        })
        const json = await res.json()
        const response: GetTagResponse = json
        response.tag.time = new Date(response.tag.time)
        return response
    }
    public async get_text(get_text_request: GetTextRequest): Promise<GetTextResponse> {
        const res = await fetch(get_text_address, {
            method: get_text_method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(get_text_request),
        })
        const json = await res.json()
        const response: GetTextResponse = json
        response.text.time = new Date(response.text.time)
        return response
    }

    public async get_tag_names(get_tag_names_request: GetTagNamesRequest): Promise<GetTagNamesResponse> {
        const res = await fetch(get_tag_names_address, {
            method: get_tag_names_method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(get_tag_names_request),
        })
        const json = await res.json()
        const response: GetTagNamesResponse = json
        return response
    }

    public async get_board_names(get_board_names_request: GetBoardNamesRequest): Promise<GetBoardNamesResponse> {
        const res = await fetch(get_board_names_address, {
            method: get_board_names_method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(get_board_names_request),
        })
        const json = await res.json()
        const response: GetBoardNamesResponse = json
        return response
    }

    public async get_application_config(get_application_config_request: GetApplicationConfigRequest): Promise<GetApplicationConfigResponse> {
        const res = await fetch(get_application_config_address, {
            method: get_application_config_method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(get_application_config_request),
        })
        const json = await res.json()
        const response: GetApplicationConfigResponse = json
        return response
    }
}