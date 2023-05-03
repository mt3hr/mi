<template>
    <textVue :text="t" />
    <tagVue :tag="tag" />
    <add_tag_dialog :show="true" :task_info="task_info" />
    <tag_list :checked_tags="[]" :option="option" />
    <board_list :option="option" />
    <board_task :task_info="task_info" />
    <detail_task :task_info="task_info" />
    <add_task_dialog ref="add_task_dialog_ref" />
    <board :board_info="board_info" :task_infos="task_infos" />
    <sort_condition_selectbox @updated_sort_type="(sort_type) => { consolelog(sort_type) }" @errors="() => { }" />
</template>

<script setup lang="ts">
import sort_condition_selectbox from './sidebar/sort_condition_selectbox.vue';
import Tag from '@/api/data_struct/Tag';
import textVue from './text/text.vue';
import tagVue from './tag/tag.vue';
import Text from '@/api/data_struct/Text';
import { Ref, ref, nextTick } from 'vue';
import TaskInfo from '@/api/data_struct/TaskInfo';
import add_tag_dialog from './dialog/add_tag_dialog.vue';
import tag_list from './sidebar/tag_list.vue';
import board_list from './sidebar/board_list.vue'
import ApplicationConfig from '@/api/data_struct/ApplicationConfig';
import board_task from './task/board_task.vue';
import add_task_dialog from './dialog/add_task_dialog.vue';
import detail_task from './task/detail_task.vue';
import board from './board/board.vue';
import BoardInfo from '@/api/data_struct/BoardInfo';

let t: Ref<Text> = ref(new Text())
t.value.text = "hoge"
let tag: Ref<Tag> = ref(new Tag())
tag.value.tag = "tag"
let task_info: Ref<TaskInfo> = ref(new TaskInfo())
let new_task_info = new TaskInfo()
new_task_info.task_title_info.title = "わかる"
new_task_info.limit_info.limit = new Date(0)
task_info.value = new_task_info
const add_task_dialog_ref = ref<InstanceType<typeof add_task_dialog> | null>(null);
nextTick(() => {
    add_task_dialog_ref.value?.show()
})

const board_info: Ref<BoardInfo> = ref(new BoardInfo())
board_info.value.board_name = "Inbox"
const task_infos: Ref<Array<TaskInfo>> = ref(new Array<TaskInfo>())
task_infos.value.push(task_info.value)


//TODO タスクのタイトルが更新されないんだが？

let option: Ref<ApplicationConfig> = ref(new ApplicationConfig())
let tag_struct_object: Ref<any> = ref({
    "hoge": "tag",
    "fuga": {
        "piyo": "tag"
    }
})
let board_struct_object: Ref<any> = ref({
    "hoge": "board",
    "fuga": {
        "piyo": "board"
    }
})
option.value.tag_struct = tag_struct_object.value
option.value.board_struct = board_struct_object
function consolelog(value: any) {
    console.log(value)
}
</script>

<style></style>