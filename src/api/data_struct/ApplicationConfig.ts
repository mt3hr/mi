export default class ApplicationConfig {
    public hidden_tags: Array<string> = new Array<string>()
    public un_check_tags: Array<string> = new Array<string>()
    public default_board_name: string = ""
    public board_struct: any = null
    public tag_struct: any = null
    public enable_hot_reload: boolean = false
}