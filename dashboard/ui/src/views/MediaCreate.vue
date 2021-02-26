<template>
  <v-row>
    <v-col cols="12" md="12">
      <MediaForm />
    </v-col>
  </v-row>
</template>

<script lang="ts">
import Vue from "vue";
import MediaForm from "../components/MediaForm.vue";

export default Vue.extend({
  components: { MediaForm }
  // data: () => ({}),
  // methods: {},
  // async created() {}
});
</script>
