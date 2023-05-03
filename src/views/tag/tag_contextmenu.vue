<template>
    <v-menu :style="style" v-model="is_show">
        <v-list>
            <v-list-item @click="delete_tag">
                <v-list-item-title>削除</v-list-item-title>
            </v-list-item>
        </v-list>
    </v-menu>
</template>

<script setup lang="ts">
import DeleteTagRequest from '@/api/DeleteTagRequest';
import MiServerAPI from '@/api/MiServerAPI';
import Tag from '@/api/data_struct/Tag';
import { Ref, ref, watch } from 'vue';

interface Props {
    tag: Tag
    x: number
    y: number
}

const props = defineProps<Props>()
const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'deleted_tag'): void
}>()

let style: Ref<string> = ref(generate_style())
let is_show: Ref<boolean> = ref(false)
defineExpose({show})

watch(() => props.x, ()=> {
    style.value = generate_style()
})
watch(() => props.y, ()=> {
    style.value = generate_style()
})

function show() {
    is_show.value = true
}
function generate_style(): string {
    return `{ position: absolute; left: ${props.x}px; top: ${props.y}px; }`
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
            emits("deleted_tag")
        })
}

function emit_errors(errors: Array<string>) {
    emits("errors", errors)
}
</script>

<style>
</style>