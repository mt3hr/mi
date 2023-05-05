<template>
    <v-card @contextmenu.prevent.stop="show_contextmenu" class="detail_task_card pa-0 ma-0">
        <v-card-title>
            タスク詳細
        </v-card-title>
        <table>
            <tr>
                <td v-for="tag_data in tags" :key="tag_data.id">
                    <tag :tag="tag_data" @deleted_tag="deleted_tag" @errors="emit_errors" />
                </td>
            </tr>
        </table>
        <table>
            <tr>
                <td>
                    <v-checkbox v-model="check" class="pa-0 ma-0" />
                </td>
                <td class="pa-0 ma-0">
                    {{ title }}
                </td>
                <td v-if="limit">
                    <p>期限: {{ limit.toLocaleString() }}</p>
                </td>
            </tr>
        </table>

        <table>
            <tr v-for="text_data in texts" :key="text_data.id">
                <td>
                    <text_view :text="text_data" @deleted_text="deleted_text" @errors="emit_errors" />
                </td>
            </tr>
        </table>
        <task_contextmenu :task_info="task_info" :x="x_contextmenu" :y="y_contextmenu" @added_tag="emit_added_tag"
            @added_text="emit_added_text" @copied_task_id="emit_copied_task_id" @updated_task="emit_updated_task"
            @deleted_task="emit_deleted_task" @errors="emit_errors" ref="task_context_menu_ref" />
    </v-card>
</template>

<script setup lang="ts">
import CheckStateInfo from '@/api/data_struct/CheckStateInfo';
import TaskInfo from '@/api/data_struct/TaskInfo'
import MiServerAPI from '@/api/MiServerAPI';
import UpdateTaskRequest from '@/api/UpdateTaskRequest';
import { Ref, ref, watch, nextTick, defineExpose } from 'vue';
import task_contextmenu from './task_contextmenu.vue';
import generate_uuid from '@/generate_uuid';
import tag from '../tag/tag.vue';
import text_view from '../text/text.vue'
import Tag from '@/api/data_struct/Tag';
import Text from '@/api/data_struct/Text';
import GetTagsRelatedTaskRequest from '@/api/GetTagsRelatedTaskRequest';
import GetTextsRelatedTaskRequest from '@/api/GetTextsRelatedTaskRequest';

interface Props {
    task_info: TaskInfo
}

const props = defineProps<Props>()
const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'copied_task_id', task_info: TaskInfo): void
    (e: 'added_tag'): void
    (e: 'added_text'): void
    (e: 'updated_task', old_task_info: TaskInfo, new_task_info: TaskInfo): void
    (e: 'deleted_task', task_info: TaskInfo): void
    (e: 'deleted_tag'): void
    (e: 'deleted_text'): void
}>()

defineExpose({
    update_tags,
    update_texts
})

let old_task_info: Ref<TaskInfo> = ref(props.task_info)
let check: Ref<boolean> = ref(props.task_info.check_state_info.is_checked)
let title: Ref<string> = ref(props.task_info.task_title_info.title)
let limit: Ref<Date | null> = ref(props.task_info.limit_info.limit)
let tags: Ref<Array<Tag>> = ref(new Array<Tag>())
let texts: Ref<Array<Text>> = ref(new Array<Text>())
let x_contextmenu: Ref<number> = ref(0)
let y_contextmenu: Ref<number> = ref(0)
const task_context_menu_ref = ref<InstanceType<typeof task_contextmenu> | null>(null);
let updating_task_info = false
update_tags()
update_texts()

watch(() => props.task_info, () => {
    updating_task_info = true
    check.value = props.task_info.check_state_info.is_checked
    title.value = props.task_info.task_title_info.title
    limit.value = props.task_info.limit_info.limit
    update_tags()
    update_texts()
    old_task_info.value = props.task_info
    nextTick(() => updating_task_info = false)
})

watch(check, () => {
    if (updating_task_info) return
    const api = new MiServerAPI()

    const new_task_info = new TaskInfo()
    new_task_info.task = props.task_info.task
    new_task_info.task_title_info = props.task_info.task_title_info
    new_task_info.check_state_info = new CheckStateInfo()
    new_task_info.limit_info = props.task_info.limit_info
    new_task_info.board_info = props.task_info.board_info

    new_task_info.check_state_info.check_state_id = generate_uuid()
    new_task_info.check_state_info.task_id = new_task_info.task.task_id
    new_task_info.check_state_info.updated_time = new Date(Date.now())
    new_task_info.check_state_info.is_checked = check.value

    const update_request = new UpdateTaskRequest()
    update_request.task_info = new_task_info

    api.update_task(update_request)
        .then(res => {
            if (res.errors && res.errors.length != 0) {
                emit_errors(res.errors)
                return
            }
            emit_updated_task(props.task_info, new_task_info)
        })

})
function update_tags() {
    const api = new MiServerAPI()
    const request = new GetTagsRelatedTaskRequest()
    request.task_id = props.task_info.task.task_id
    api.get_tags_related_task(request)
        .then(res => {
            if (res.errors && res.errors.length != 0) {
                emit_errors(res.errors)
                return
            }
            tags.value = res.tags
        })
}
function update_texts() {
    const api = new MiServerAPI()
    const request = new GetTextsRelatedTaskRequest()
    request.task_id = props.task_info.task.task_id
    api.get_texts_related_task(request)
        .then(res => {
            if (res.errors && res.errors.length != 0) {
                emit_errors(res.errors)
                return
            }
            texts.value = res.texts
        })
}
function deleted_tag() {
    update_tags()
    emit_deleted_tag()
}
function deleted_text() {
    update_texts()
    emit_deleted_text()
}
function show_contextmenu(e: MouseEvent) {
    x_contextmenu.value = e.x
    y_contextmenu.value = e.y
    task_context_menu_ref.value!.show()
}
function emit_errors(errors: Array<string>) {
    emits("errors", errors)
}
function emit_copied_task_id() {
    emits("copied_task_id", props.task_info)
}
function emit_added_tag() {
    emits("added_tag")
}
function emit_added_text() {
    emits("added_text")
}
function emit_updated_task(old_task_info: TaskInfo, new_task_info: TaskInfo) {
    emits("updated_task", old_task_info, new_task_info)
}
function emit_deleted_task(deleted_task_info: TaskInfo) {
    emits("deleted_task", deleted_task_info)
}
function emit_deleted_text() {
    emits("deleted_text")
}
function emit_deleted_tag() {
    emits("deleted_tag")
}
</script>

<style>
.v-checkbox>.v-input__details {
    height: 0 !important;
    max-height: 0 !important;
    min-height: 0 !important;
}
</style>