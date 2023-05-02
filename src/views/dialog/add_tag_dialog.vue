<template>
    <v-dialog v-model="is_show">
        <v-card class="pa-5">
            <v-card-title>
                タグ追加
            </v-card-title>
            <v-text-field v-model="tag_name" @keypress.enter="submit" tabindex="101" />
            <v-card-actions>
                <v-row>
                    <v-col cols="auto">
                        <v-btn @click="close_dialog" tabindex="103">
                            キャンセル
                        </v-btn>
                    </v-col>
                    <v-spacer />
                    <v-col cols="auto">
                        <v-btn @click="submit" tabindex="102">
                            追加
                        </v-btn>
                    </v-col>
                </v-row>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import AddTagRequest from '@/api/AddTagRequest';
import MiServerAPI from '@/api/MiServerAPI';
import TaskInfo from '@/api/data_struct/TaskInfo';
import { Ref, ref, watch } from 'vue';

interface Props {
    task_info: TaskInfo
}

const props = defineProps<Props>()
const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'added_tag'): void
}>()

let tag_name: Ref<string> = ref("")
let is_show: Ref<boolean> = ref(false)

defineExpose({show})

watch(() => is_show.value, () => {
    is_show.value = is_show.value
})

function close_dialog() {
    is_show.value = false
}
function submit() {
    if (tag_name.value == "") {
        return
    }
    const api = new MiServerAPI()
    const request = new AddTagRequest()
    request.tag = tag_name.value
    request.task_id = props.task_info.task.task_id
    api.add_tag(request)
        .then(res => {
            if (res.errors.length != 0) {
                emit_errors(res.errors)
                return
            }
            emit_added_tag()
            clear_fields()
            close_dialog()
        })
}
function show() {
    is_show.value = true
}
function clear_fields() {
    tag_name.value = ""
}

function emit_errors(errors: Array<string>) {
    emits("errors", errors)
}
function emit_added_tag() {
    emits("added_tag")
}
</script>

<style></style>