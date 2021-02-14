<template>
    <v-row>
      <v-col cols="12" md="6" sm="12">
        <v-card height="100%" outlined>
          <v-card-title>Welcome to Gilfoyle Dashboard!</v-card-title>
          <v-card-text
            >Here you can manage your content and view metrics about your
            instance. At any time, if needed, you can check out the
            <a :href="references.docsURL" target="_blank">documentation</a>.
          </v-card-text>
        </v-card>
      </v-col>

      <v-col cols="12" md="6" sm="12">
        <HealthGraph></HealthGraph>
      </v-col>

      <v-col md="12">
        <MediaList></MediaList>
      </v-col>
    </v-row>
</template>

<script lang="ts">
import { Vue } from "vue-property-decorator";
import MediaList from "../components/MediaList.vue";
import HealthGraph from "../components/HealthGraph.vue";
import config from "../config";

interface Data {
  references: any;
}

export default Vue.extend({
  components: { MediaList, HealthGraph },
  data: (): Data => ({
    references: config.links
  })
});
</script>
