import MiResponse from "./MiResponse";
import ApplicationConfig from "./data_struct/ApplicationConfig";

export default class GetApplicationConfigResponse extends MiResponse {
    public application_config: ApplicationConfig = new ApplicationConfig()
}