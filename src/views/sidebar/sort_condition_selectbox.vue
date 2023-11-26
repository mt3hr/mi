<template>
    <v-select v-model="sort_type" :label="'ソート'" :items="sort_types" item-title="title" item-value="value"
        @update:model-value="emit_updated_sort_type" />
</template>

<script setup lang="ts">
import SortType from '@/api/data_struct/SortType';
import { Ref, ref, nextTick, watch } from 'vue';
class SortTypeSelectModel {
    title: string = ""
    value: SortType = SortType.CreatedTimeDesc
}
const sort_type_select_model = new Array<SortTypeSelectModel>()
const sort_type_created_time_desc = new SortTypeSelectModel()
sort_type_created_time_desc.title = "作成日時順"
sort_type_created_time_desc.value = SortType.CreatedTimeDesc
const sort_type_limit_time_asc = new SortTypeSelectModel()
sort_type_limit_time_asc.title = "期限順"
sort_type_limit_time_asc.value = SortType.LimitTimeAsc
const sort_type_start_time_desc = new SortTypeSelectModel()
sort_type_start_time_desc.title = "開始日時順"
sort_type_start_time_desc.value = SortType.StartTimeDesc
const sort_type_end_time_desc = new SortTypeSelectModel()
sort_type_end_time_desc.title = "終了日時順"
sort_type_end_time_desc.value = SortType.EndTimeDesc
sort_type_select_model.push(sort_type_created_time_desc)
sort_type_select_model.push(sort_type_limit_time_asc)
sort_type_select_model.push(sort_type_start_time_desc)
sort_type_select_model.push(sort_type_end_time_desc)

const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'updated_sort_type', sort_type: SortType): void
}>()

const sort_types: Ref<Array<SortTypeSelectModel>> = ref(sort_type_select_model)
const sort_type: Ref<SortType> = ref(SortType.CreatedTimeDesc)

defineExpose({
    set_sort_type_by_application,
    get_sort_type
})

function get_sort_type(): SortType {
    return sort_type.value
}
function set_sort_type_by_application(new_sort_type: SortType): void {
    sort_type.value = new_sort_type
}

function emit_errors(errors: Array<string>) {
    emits("errors", errors)
}
function emit_updated_sort_type() {
    emits("updated_sort_type", sort_type.value)
}
</script>

<style></style>