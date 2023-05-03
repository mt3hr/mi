<template>
    <v-select v-model="check_state" :label="'チェック状態'" :items="check_states" item-title="title" item-value="value"
        @update:model-value="emit_updated_check_condition" />
</template>

<script setup lang="ts">
import CheckState from '@/api/data_struct/CheckState';
import { Ref, ref, nextTick, watch } from 'vue';
class CheckStateSelectModel {
    title: string = ""
    value: CheckState = CheckState.NoCheckOnly
}
const check_state_select_model = new Array<CheckStateSelectModel>()
const check_state_no_check_only = new CheckStateSelectModel()
check_state_no_check_only.title = "未チェックのみ"
check_state_no_check_only.value = CheckState.NoCheckOnly
const check_state_check_only = new CheckStateSelectModel()
check_state_check_only.title = "チェック済みのみ"
check_state_check_only.value = CheckState.CheckOnly
const check_state_all = new CheckStateSelectModel()
check_state_all.title = "全て"
check_state_all.value = CheckState.All
check_state_select_model.push(check_state_no_check_only)
check_state_select_model.push(check_state_check_only)
check_state_select_model.push(check_state_all)

const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'updated_check_condition', check_state: CheckState): void
}>()

const check_states: Ref<Array<CheckStateSelectModel>> = ref(check_state_select_model)
const check_state: Ref<CheckState> = ref(CheckState.NoCheckOnly)

defineExpose({
    set_check_state_by_application, 
    get_check_state
})

function get_check_state(): CheckState {
    return check_state.value
}
function set_check_state_by_application(new_check_state: CheckState): void {
    check_state.value = new_check_state
}

function emit_errors(errors: Array<string>) {
    emits("errors", errors)
}
function emit_updated_check_condition() {
    emits("updated_check_condition", check_state.value)
}
</script>

<style></style>