<template>
    <textVue :text="t" />
    <tagVue :tag="tag" />
    <add_tag_dialog :show="true" :task_info="task_info"></add_tag_dialog>
    <table>
    <tag_struct :struct="tag_struct_object" :group_name="''" :open="true"></tag_struct>
    </table>
</template>

<script setup lang="ts">
import Tag from '@/api/data_struct/Tag';
import textVue from './text/text.vue';
import tagVue from './tag/tag.vue';
import Text from '@/api/data_struct/Text';
import { Ref, ref } from 'vue';
import TaskInfo from '@/api/data_struct/TaskInfo';
import add_tag_dialog from './dialog/add_tag_dialog.vue';
import tag_struct from './sidebar/tag_struct.vue';

let t: Ref<Text> = ref(new Text())
t.value.text = "hoge"
let tag: Ref<Tag> = ref(new Tag())
tag.value.tag = "tag"
let task_info: Ref<TaskInfo> = ref(new TaskInfo())

let tag_struct_object: Ref<any> = ref({
    "hoge": "tag",
    "fuga": {
        "piyo": "tag"
    }
})
</script>

<style></style>