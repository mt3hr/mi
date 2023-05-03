<template>
    <h2>タグ</h2>
    <table class="taglist">
        <tag_struct ref="tag_struct_ref" :group_name="''" :struct="tag_structure" :open="true"
            @updated_check_items_by_user="updated_checked_tags" @click_items_by_user="check_only_tags" />
    </table>
</template>
<script setup lang="ts">
import { Ref, ref, watch, nextTick } from 'vue';
import MiServerAPI from '@/api/MiServerAPI';
import tag_struct from './tag_struct.vue';
import GetTagNamesRequest from '@/api/GetTagNamesRequest';
import ApplicationConfig from '@/api/data_struct/ApplicationConfig';

interface Props {
    option: ApplicationConfig
}

const props = defineProps<Props>()
const emits = defineEmits<{
    (e: 'errors', errors: Array<string>): void
    (e: 'updated_by_user'): void
    (e: 'updated_checked_tags', tags: Array<string>): void
}>()

let tags: Ref<any> = ref({})
let tag_structure: Ref<any> = ref({})
let check_all: Ref<boolean> = ref(true)
const tag_struct_ref = ref<InstanceType<typeof tag_struct> | null>(null);
const checked_tags: Ref<Array<string>> = ref(new Array<string>());

defineExpose({
    set_checked_tags_by_application,
    get_checked_tags
})


nextTick(() => {
    update_tags_promise()
        .then(() => { return check_all_tags_promise() })
        .then(() => { return update_tag_struct_promise() })
        .then(() => emits('updated_by_user'))
})

watch(() => checked_tags, () => {
    for (let i = 0; i < tags.value.length; i++) {
        let tag: any = tags.value[i]
        tag.check = false
    }
    for (let i = 0; i < tags.value.length; i++) {
        let tag: any = tags.value[i]
        for (let j = 0; j < checked_tags.value.length; j++) {
            let checked_tag = checked_tags.value[j]
            if (tag.tag == checked_tag) {
                tag.check = true
            }
        }
    }
    update_tag_struct_promise()
})


// 渡されたタグのチェック状態を更新します。
function updated_checked_tags(checked_tags: Array<any>, check: boolean) {
    for (let j = 0; j < checked_tags.length; j++) {
        let tag = checked_tags[j]
        for (let i = 0; i < tags.value.length; i++) {
            if (tags.value[i].tag === tag) {
                tags.value[i].check = check
                break
            }
        }
    }
    update_tag_struct_promise()
        .then(() => emits('updated_by_user'))
}
// 全てのtagにチェックを入れる
function check_all_tags_promise() {
    let f = (option: any) => {
        for (let i = 0; i < tags.value.length; i++) {
            let tag = tags.value[i]
            let check = true
            if (option.un_check_tags) {
                for (let i = 0; i < option.un_check_tags.length; i++) {
                    let uncheck_tag = option.un_check_tags[i]
                    if (uncheck_tag == tag.tag) {
                        check = false
                        break
                    }
                }
            }
            tag.check = check
        }
    }
    return new Promise(resolve => { return resolve(null) })
        .then(() => { if (tags.value.length === 0) return update_tags_promise() })
        .then(() => f(props.option))
        .then(() => { return update_tag_struct_promise() })
        .then(() => notify_checked_tags())
}
// tag_structをkv_taglist_tagsの取り扱える形に変換し、更新します。
function update_tag_struct_promise() {
    return new Promise(resolve => { return resolve(null) })
        .then(() => {
            // tag structを変換してから収めます。
            /*
            {
              "no tag": "tag",
              "log": {
                "の": "tag",
                "ぢ": "tag",
              }
            }
            から
            {
              no tag: {tag: "no tag", check: true},
              log: {
                  "の": {tag: "の", check: true},
                  "ぢ": {tag: "ぢ", check: true},
              },
            }
            に。
            */
            let structed_tags: any = []
            let f = (tag_or_parent: any, tag_name: string) => { }
            let func = (tag_or_parent: any, tagname: string) => {
                if (tag_or_parent === 'tag') {
                    let check = false
                    for (let i = 0; i < tags.value.length; i++) {
                        if (tags.value[i].tag == tagname) {
                            check = tags.value[i].check
                            break
                        }
                    }
                    structed_tags.push(tagname)
                    return {
                        key: tagname,
                        indeterminate: false,
                        check: check,
                    }
                } else {
                    let tag_struct: any = {}
                    Object.keys(tag_or_parent).forEach(parent => {
                        tag_struct[parent] = f(tag_or_parent[parent], parent)
                    })
                    return tag_struct
                }
            }
            f = func
            let tag_struct_
            if (props.option.tag_struct) {
                tag_struct_ = f(props.option.tag_struct, "")
            } else {
                tag_struct_ = f({}, "")
            }

            tags.value.forEach((tag: any) => {
                let exist = false
                for (let i = 0; i < structed_tags.length; i++) {
                    if (tag.tag == structed_tags[i]) {
                        exist = true
                        break
                    }
                }
                if (!exist) {
                    tag_struct_[tag.tag] = {
                        key: tag.tag,
                        indeterminate: false,
                        check: tag.check,
                    }
                }

            })
            structed_tags.forEach((tag: any) => {
                let exist = false
                for (let i = 0; i < tags.value.length; i++) {
                    if (tag == tags.value[i].tag) {
                        exist = true
                        break
                    }
                }
                if (!exist) {
                    let check = true
                    for (let i = 0; i < props.option.un_check_tags.length; i++) {
                        if (tag === props.option.un_check_tags[i]) {
                            check = false
                            break
                        }
                    }
                    tags.value.push({
                        tag: tag,
                        check: check,
                    })
                }
            })
            tag_structure.value = tag_struct_
        }).then(() => notify_checked_tags())
}
// タグを最新の状態に更新します。
// タグの選択はすべてfalseに初期化されます。
function update_tags_promise() {
    let api = new MiServerAPI()
    return api.get_tag_names(new GetTagNamesRequest())
        .then((res) => {
            let tagsTemp: any = []
            res.tag_names.forEach((tag: any) => {
                let t = {
                    tag: tag,
                    check: true
                }
                tagsTemp.push(t)
            })
            tags.value = tagsTemp
        })
        .then(() => { return update_tag_struct_promise() })
        .then(() => notify_checked_tags())
}
// 現在選択されているタグを[]stringで取得します。
// notify_checked_tagsから呼び出されます。
function get_selected_tags(): Array<string> {
    return tag_struct_ref.value!.get_selected_items()
}
// 現在選択されているタグをemitで通知します。
// このメソッドはチェックボックスの状態更新時そのたびにテンプレートから呼び出されます。
function notify_checked_tags() {
    const tags = get_selected_tags()
    emits('updated_checked_tags', tags)
}
// 渡された全てのtagに飲みチェックを入れ、他のタグのチェックを外します。
function check_only_tags(tags_: any) {
    for (let i = 0; i < tags.value.length; i++) {
        tags.value[i].check = false
    }
    for (let i = 0; i < tags_.length; i++) {
        let tag = tags_[i]
        for (let j = 0; j < tags.value.length; j++) {
            let t = tags.value[j]
            if (t.tag === tag) {
                t.check = true
                continue
            }
        }
    }
    update_tag_struct_promise()
        .then(() => emits('updated_by_user'))
}
function get_checked_tags(): Array<string> {
    return checked_tags.value
}
function set_checked_tags_by_application(new_checked_tags: Array<string>): void {
    checked_tags.value = new_checked_tags
}
</script>

<style></style>