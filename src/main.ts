// Plugins
import { registerPlugins } from "@/plugins";
import { applyMetadata } from "@/utils/metadata";

// Components
import App from "./App.vue";

// Composables
import { createApp } from "vue";

// Styles
import "./style.css";

const app = createApp(App);

applyMetadata();

registerPlugins(app);

app.mount("#app");
