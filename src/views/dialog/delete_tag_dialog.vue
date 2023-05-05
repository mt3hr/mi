<template>
    <v-dialog v-model="is_show" :width="500">
        <v-card class="pa-5">
            <v-card-title>
                タグ削除
            </v-card-title>
            <v-card-text>{{ tag.tag }}</v-card-text>
            <v-card-actions>
                <v-row>
                    <v-col cols="auto">
                        <v-btn @click="close_dialog" tabindex="103">
                            キャンセル
                        </v-btn>
                    </v-col>
                    <v-spacer />
                    <v-col cols="auto">
                        <v-btn @click="delete_tag" tabindex="102">
                            削除
                        </v-btn>
                    </v-col>
                </v-row>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import DeleteTagRequest from '@/api/DeleteTagRequest';
import MiServerAPI from '@/api/MiServerAPI';
import Tag from '@/api/data_struct/Tag';
import TaskInfo from '@/api/data_struct/TaskInfo';
import { Ref, ref, watch } from 'vue';

interface Props {
    tag: Tag
    task_info: TaskInfo
}

const props = defineProps<Props>()
const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'deleted_tag'): void
}>()

let is_show: Ref<boolean> = ref(false)

defineExpose({show})

watch(() => is_show.value, () => {
    is_show.value = is_show.value
})

function close_dialog() {
    is_show.value = false
}
function delete_tag() {
    const api = new MiServerAPI()
    const request = new DeleteTagRequest()
    request.tag_id = props.tag.id
    api.delete_tag(request)
        .then(res => {
            if (res.errors && res.errors.length != 0) {
                emit_errors(res.errors)
                return
            }
            emit_deleted_tag()
            close_dialog()
        })
}
function show() {
    is_show.value = true
}

function emit_errors(errors: Array<string>) {
    emits("errors", errors)
}
function emit_deleted_tag() {
    emits("deleted_tag")
}
</script>

<style>
</style>