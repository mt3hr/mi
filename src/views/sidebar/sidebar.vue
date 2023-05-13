<template>
    <div>
        <word_search_textbox @errors="emit_errors" @updated_search_word="updated_search_word"
            ref="word_search_textbox_ref" />
        <check_condition_selectbox @errors="emit_errors" @updated_check_condition="updated_check_condition"
            ref="check_condition_selectbox_ref" />
        <sort_condition_selectbox @errors="emit_errors" @updated_sort_type="updated_sort_type"
            ref="sort_condition_selectbox_ref" />
        <board_list v-if="option" :option="option" @errors="emit_errors" @updated_by_user="updated_boards_by_user"
            @clicked_board="clicked_board" ref="board_list_ref" />
        <tag_list v-if="option" :option="option" @errors="emit_errors" @updated_by_user="updated_tags_by_user"
            @updated_checked_tags="updated_checked_tags" ref="tag_list_ref" />
    </div>
</template>

<script setup lang="ts">
import { Ref, ref, watch, nextTick } from 'vue';
import check_condition_selectbox from './check_condition_selectbox.vue';
import sort_condition_selectbox from './sort_condition_selectbox.vue';
import word_search_textbox from './word_search_textbox.vue';
import board_list from './board_list.vue';
import tag_list from './tag_list.vue';
import CheckState from '@/api/data_struct/CheckState';
import SortType from '@/api/data_struct/SortType';
import ApplicationConfig from '@/api/data_struct/ApplicationConfig';
import TaskSearchQuery from '@/api/data_struct/TaskSearchQuery';

interface Props {
    option: ApplicationConfig
}
const props = defineProps<Props>()
const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'updated_check_condition', check_state: CheckState): void
    (e: 'updated_sort_type', sort_type: SortType): void
    (e: 'updated_search_word', word: string): void
    (e: 'updated_boards_by_user'): void
    (e: 'clicked_board', board: any): void
    (e: 'updated_tags_by_user'): void
    (e: 'updated_checked_tags', tags: Array<string>): void
}>()

const check_condition_selectbox_ref = ref<InstanceType<typeof check_condition_selectbox> | null>(null);
const sort_condition_selectbox_ref = ref<InstanceType<typeof sort_condition_selectbox> | null>(null);
const word_search_textbox_ref = ref<InstanceType<typeof word_search_textbox> | null>(null);
const board_list_ref = ref<InstanceType<typeof board_list> | null>(null);
const tag_list_ref = ref<InstanceType<typeof tag_list> | null>(null);

defineExpose({
    set_search_word_by_application,
    get_search_word,
    set_sort_type_by_application,
    get_sort_type,
    set_check_state_by_application,
    get_check_state,
    set_checked_tags_by_application,
    get_checked_tags,
    construct_task_search_query,
    check_all_tags,
    update_tag_struct_promise,
    update_board_struct_promise
})

function updated_check_condition(updated_check_state: CheckState) {
    emit_updated_check_condition(updated_check_state)
}
function updated_sort_type(updated_sort_type: SortType) {
    emit_updated_sort_type(updated_sort_type)
}
function updated_search_word(updated_word: string) {
    emit_updated_search_word(updated_word)
}
function updated_boards_by_user() {
    emit_updated_boards_by_user()
}
function clicked_board(updated_board: string) {
    emit_clicked_board(updated_board)
}
function updated_tags_by_user() {
    emit_updated_tags_by_user()
}
function updated_checked_tags(checked_tags: Array<string>) {
    emit_updated_checked_tags(checked_tags)
}
function construct_task_search_query(): TaskSearchQuery {
    const query = new TaskSearchQuery()
    query.check_state = check_condition_selectbox_ref.value?.get_check_state()!
    query.sort_type = sort_condition_selectbox_ref.value?.get_sort_type()!
    query.tags = tag_list_ref.value?.get_checked_tags()!
    query.word = word_search_textbox_ref.value?.get_search_word()!
    return query
}

function get_search_word(): string {
    return word_search_textbox_ref.value?.get_search_word()!
}
function set_search_word_by_application(new_search_word: string): void {
    word_search_textbox_ref.value?.set_search_word_by_application(new_search_word)
}
function get_sort_type(): SortType {
    return sort_condition_selectbox_ref.value?.get_sort_type()!
}
function set_sort_type_by_application(new_sort_type: SortType): void {
    sort_condition_selectbox_ref.value?.set_sort_type_by_application(new_sort_type)
}
function get_check_state(): CheckState {
    return check_condition_selectbox_ref.value?.get_check_state()!
}
function set_check_state_by_application(new_check_state: CheckState): void {
    check_condition_selectbox_ref.value?.set_check_state_by_application(new_check_state)
}
function get_checked_tags(): Array<string> {
    return tag_list_ref.value?.get_checked_tags()!
}
function set_checked_tags_by_application(new_checked_tags: Array<string>) {
    tag_list_ref.value?.set_checked_tags_by_application(new_checked_tags)
}

function emit_errors(errors: Array<string>) {
    emits("errors", errors)
}
function emit_updated_check_condition(updated_check_state: CheckState) {
    emits("updated_check_condition", updated_check_state)
}
function emit_updated_sort_type(updated_sort_type: SortType) {
    emits("updated_sort_type", updated_sort_type)
}
function emit_updated_search_word(updated_word: string) {
    emits("updated_search_word", updated_word)
}
function emit_updated_boards_by_user() {
    emits("updated_boards_by_user")
}
function emit_clicked_board(updated_board: string) {
    emits("clicked_board", updated_board)
}
function emit_updated_tags_by_user() {
    emits("updated_tags_by_user")
}
function emit_updated_checked_tags(updated_checked_tags: Array<string>) {
    emits("updated_checked_tags", updated_checked_tags)
}
async function check_all_tags() {
    await tag_list_ref.value?.check_all_tags()
}
async function update_tag_struct_promise() {
    const checked_tags = get_checked_tags()
    await tag_list_ref.value?.update_tags_promise()
    set_checked_tags_by_application(checked_tags)
}
async function update_board_struct_promise() {
    await board_list_ref.value?.update_boards_promise()
}
</script>

<style></style>