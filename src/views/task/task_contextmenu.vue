<template>
    <v-menu :style="style" v-model="is_show">
        <v-list>
            <v-list-item @click="show_add_tag_dialog">
                <v-list-item-title>タグ追加</v-list-item-title>
            </v-list-item>
            <v-list-item @click="show_add_text_dialog">
                <v-list-item-title>テキスト追加</v-list-item-title>
            </v-list-item>
            <v-list-item @click="copy_task_id_to_clipboard">
                <v-list-item-title>IDをコピー</v-list-item-title>
            </v-list-item>
            <v-list-item @click="show_edit_task_dialog">
                <v-list-item-title>編集</v-list-item-title>
            </v-list-item>
            <v-list-item @click="show_delete_task_dialog">
                <v-list-item-title>削除</v-list-item-title>
            </v-list-item>
        </v-list>
    </v-menu>
    <add_tag_dialog :task_info="task_info" ref="add_tag_dialog_ref" @errors="emit_errors" @added_tag="emit_added_tag" />
    <add_text_dialog :task_info="task_info" ref="add_text_dialog_ref" @errors="emit_errors" @added_text="emit_added_text" />
    <edit_task_dialog :task_info="task_info" ref="edit_task_dialog_ref" @errors="emit_errors"
        @updated_task="emit_updated_task" />
    <delete_task_dialog :task_info="task_info" ref="delete_task_dialog_ref" @errors="emit_errors"
        @deleted_task="emit_deleted_task" />
</template>

<script setup lang="ts">
import add_tag_dialog from '../dialog/add_tag_dialog.vue';
import add_text_dialog from '../dialog/add_text_dialog.vue';
import edit_task_dialog from '../dialog/edit_task_dialog.vue';
import delete_task_dialog from '../dialog/delete_task_dialog.vue';
import { Ref, ref, watch } from 'vue';
import TaskInfo from '@/api/data_struct/TaskInfo';

interface Props {
    task_info: TaskInfo
    x: number
    y: number
}

const props = defineProps<Props>()
const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'copied_task_id', task_info: TaskInfo): void
    (e: 'added_tag'): void
    (e: 'added_text'): void
    (e: 'updated_task', task_info: TaskInfo): void
    (e: 'deleted_task', task_info: TaskInfo): void
}>()

let style: Ref<string> = ref(generate_style())
let is_show: Ref<boolean> = ref(false)
const add_tag_dialog_ref = ref<InstanceType<typeof add_tag_dialog> | null>(null);
const add_text_dialog_ref = ref<InstanceType<typeof add_text_dialog> | null>(null);
const edit_task_dialog_ref = ref<InstanceType<typeof edit_task_dialog> | null>(null);
const delete_task_dialog_ref = ref<InstanceType<typeof delete_task_dialog> | null>(null);

defineExpose({ show })

watch(() => props.x, () => {
    style.value = generate_style()
})
watch(() => props.y, () => {
    style.value = generate_style()
})

function show() {
    is_show.value = true
}
function generate_style(): string {
    return `{ position: absolute; left: ${props.x}px; top: ${props.y}px; }`
}
function show_add_tag_dialog() {
    add_tag_dialog_ref.value!.show()
}
function show_add_text_dialog() {
    add_text_dialog_ref.value!.show()
}
function copy_task_id_to_clipboard() {
    navigator.clipboard.writeText(props.task_info.task.task_id)
}
function show_edit_task_dialog() {
    edit_task_dialog_ref.value!.show()
    emit_copied_task_id()
}
function show_delete_task_dialog() {
    delete_task_dialog_ref.value!.show()
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
function emit_updated_task(updated_task_info: TaskInfo) {
    emits("updated_task", updated_task_info)
}
function emit_deleted_task(deleted_task_info: TaskInfo) {
    emits("deleted_task", deleted_task_info)
}
</script>

<style></style>