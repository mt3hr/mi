<template>
    <v-dialog v-model="show_new_task_dialog" class="new_task_dialog">
        <v-card class="new_task_card pa-3">
            <v-card-title>新規</v-card-title>
            <v-text-field v-model="task_title" type="text" placeholder="タイトル" />
            <v-textarea v-model="task_memo" placeholder="メモ" />
            <v-checkbox v-model="task_has_limit" label="期限" />
            <VueDatePicker v-if="task_has_limit" v-model="task_limit" placeholder="期限日時" :format="format_date_time"
                :flow="['calendar', 'time']" />
            <v-card-actions>
                <v-row>
                    <v-spacer />
                    <v-col cols="auto">
                        <v-btn>追加</v-btn>
                    </v-col>
                </v-row>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { ref, Ref, watch } from 'vue'
import VueDatePicker from '@vuepic/vue-datepicker';
import '@vuepic/vue-datepicker/dist/main.css';

let show_new_task_dialog: Ref<boolean> = ref(true)
let task_title: Ref<string> = ref("")
let task_memo: Ref<string> = ref("")
let task_has_limit: Ref<boolean> = ref(false)
let task_limit: Ref<string> = ref("")

function format_date_time(date: Date): string {
    return date.toLocaleDateString() + " " + date.toLocaleTimeString()
}

watch(task_has_limit, (new_value: boolean, old_value: boolean) => {
    if (!new_value) {
        task_limit.value = ""
    }
})
</script>

<style>
.new_task_card {
    overflow-y: hidden !important;
}

.new_task_dialog {
    max-width: 600px;
}
</style>