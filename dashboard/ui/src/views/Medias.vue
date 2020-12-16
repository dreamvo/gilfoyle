<template>
  <v-container>
    <v-row>
      <v-col md="12">
        <MediaList></MediaList>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts">
import { Vue } from "vue-property-decorator";
import MediaList from "../components/MediaList.vue";
import HealthGraph from "../components/HealthGraph.vue";

interface Data {}

export default Vue.extend({
  components: { MediaList, HealthGraph },
  data: (): Data => ({})
});
</script>
