<template>
    <!--
    <textVue :text="t" />
    <tagVue :tag="tag" />
    <add_tag_dialog :show="true" :task_info="task_info"></add_tag_dialog>
    <tag_list :checked_tags="[]" :option="option" />
    <board_list :option="option" />
    -->
    <board_task :task_info="task_info" />
</template>

<script setup lang="ts">
import Tag from '@/api/data_struct/Tag';
import textVue from './text/text.vue';
import tagVue from './tag/tag.vue';
import Text from '@/api/data_struct/Text';
import { Ref, ref } from 'vue';
import TaskInfo from '@/api/data_struct/TaskInfo';
import add_tag_dialog from './dialog/add_tag_dialog.vue';
import tag_list from './sidebar/tag_list.vue';
import board_list from './sidebar/board_list.vue'
import ApplicationConfig from '@/api/data_struct/ApplicationConfig';
import board_task from './task/board_task.vue';

let t: Ref<Text> = ref(new Text())
t.value.text = "hoge"
let tag: Ref<Tag> = ref(new Tag())
tag.value.tag = "tag"
let task_info: Ref<TaskInfo> = ref(new TaskInfo())
let new_task_info = new TaskInfo()
new_task_info.task_title_info.title = "わかる"
task_info.value = new_task_info

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
</script>

<style></style>