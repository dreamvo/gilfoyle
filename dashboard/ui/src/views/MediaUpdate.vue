<template>
  <v-row>
    <v-col cols="12" md="12">
      <MediaForm :media-id="$route.params.id" />
    </v-col>
  </v-row>
</template>

<script lang="ts">
import Vue from "vue";
import MediaForm from "../components/MediaForm.vue";

export default Vue.extend({
  components: { MediaForm }
  // data: (): Data => ({}),
  // methods: {},
  // async created() {}
});
</script>
