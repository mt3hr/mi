<template>
    <FullCalendar class="calendar" :options='calendar_options'>
        <template v-slot:eventContent='arg'>
            <p @click="emit_clicked_date(arg.event.extendedProps.date_)">{{ arg.event.title }}</p>
        </template>
    </FullCalendar>
</template>
<script lang="ts" setup>
import { Ref, ref, watch, nextTick, defineExpose } from 'vue';
import FullCalendar from '@fullcalendar/vue3'
import jaLocale from '@fullcalendar/core/locales/ja'
import dayGridPlugin from '@fullcalendar/daygrid'
import TaskInfo from '@/api/data_struct/TaskInfo';
import SortType from '@/api/data_struct/SortType';

const actual_height = window.innerHeight
const element_height = document!.querySelector('#control-height') ? document!.querySelector('#control-height')!.clientHeight : actual_height
const bar_height = (actual_height - element_height) + "px"

interface Props {
    task_infos: Array<TaskInfo>
    mode: SortType
}
const props = defineProps<Props>()
const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'clicked_date', date: Date): void
}>()

const calendar_options: Ref<any> = ref({
    plugins: [dayGridPlugin],
    initialView: "dayGridMonth",
    weekends: true,
    events: [],
    locale: jaLocale,
})

defineExpose({
    update_calendar,
})

watch(() => props.task_infos, () => {
    update_calendar()
})
watch(() => props.mode, () => {
    update_calendar()
})

function update_calendar() {
    let events: any = []
    if (props.task_infos) {
        let events_map: any = {}
        props.task_infos.forEach(taskinfo => {
            if (props.mode == SortType.CreatedTimeDesc) {
                const created_time = taskinfo.task.created_time
                if (!events_map[`${created_time.getFullYear().toString().padStart(4, '0')}-${(created_time.getMonth() + 1).toString().padStart(2, '0')}-${created_time.getDate().toString().padStart(2, '0')}`]) {
                    events_map[`${created_time.getFullYear().toString().padStart(4, '0')}-${(created_time.getMonth() + 1).toString().padStart(2, '0')}-${created_time.getDate().toString().padStart(2, '0')}`] = 0
                }
                events_map[`${created_time.getFullYear().toString().padStart(4, '0')}-${(created_time.getMonth() + 1).toString().padStart(2, '0')}-${created_time.getDate().toString().padStart(2, '0')}`] += 1
            } else if (props.mode == SortType.LimitTimeAsc) {
                const limit_time = taskinfo.limit_info.limit
                if (!limit_time) {
                    return
                }
                if (!events_map[`${limit_time.getFullYear().toString().padStart(4, '0')}-${(limit_time.getMonth() + 1).toString().padStart(2, '0')}-${limit_time.getDate().toString().padStart(2, '0')}`]) {
                    events_map[`${limit_time.getFullYear().toString().padStart(4, '0')}-${(limit_time.getMonth() + 1).toString().padStart(2, '0')}-${limit_time.getDate().toString().padStart(2, '0')}`] = 0
                }
                events_map[`${limit_time.getFullYear().toString().padStart(4, '0')}-${(limit_time.getMonth() + 1).toString().padStart(2, '0')}-${limit_time.getDate().toString().padStart(2, '0')}`] += 1
            }
        });
        Object.keys(events_map).forEach(key => {
            events.push({
                title: events_map[key],
                date: key,
                date_: key
            })
        })
    }
    calendar_options.value.events = events
}


function emit_clicked_date(datestr: string) {
    emits("clicked_date", new Date(Date.parse(`${datestr} 00:00:00`)))
}

</script>
<style>
.calendar {
    width: 370px;
    min-width: 370px;
    max-width: 370px;
    height: calc(((100vh - 50px + v-bind(bar_height)) / 2) - 10px);
    max-height: calc(((100vh - 50px + v-bind(bar_height)) / 2) - 10px);
    min-height: calc(((100vh - 50px + v-bind(bar_height)) / 2) - 10px);
}

.fc-header-toolbar.fc-toolbar {
    margin: 0 !important;
}

.fc-event.fc-daygrid-event.fc-h-event {
    max-width: 30px;
    margin-left: 10px;
    text-align: center;
}

.fc-daygrid-day-frame {
    height: 50px !important;
    max-height: 50px !important;
    min-height: 50px !important;
}

.fc-today-button.fc-button.fc-button-primary,
.fc-next-button.fc-button.fc-button-primary,
.fc-prev-button.fc-button.fc-button-primary {
    background-color: #3f51b5 !important;
}

.fc-daygrid-day-events {
    position: relative;
    top: -8px;
}

.fc-icon.fc-icon-chevron-left::before,
.fc-icon.fc-icon-chevron-right::before {
    vertical-align: top !important;
}

.fc-button:disabled {
    opacity: 100% !important;
}
</style>