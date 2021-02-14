<template>
  <v-dialog v-model="dialog" width="500">
    <template #activator="{ on, attrs }">
      <v-btn
        color="error"
        text
        class="ml-2"
        :small="btnSmall"
        v-bind="attrs"
        v-on="on"
      >
        <v-icon v-if="icon" size="16" class="mr-1">mdi-trash-can</v-icon>
        {{ btnText }}
      </v-btn>
    </template>

    <v-card outlined>
      <v-card-title class="headline lighten-2 pb-4 mb-5">
        {{ title }}
      </v-card-title>

      <v-card-text
        >Êtes vous sur de vouloir effectuer cette action ?
      </v-card-text>
      <v-card-text>Cette action est irréversible.</v-card-text>

      <v-divider></v-divider>

      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="blue-grey" text @click="dialog = false">Cancel</v-btn>
        <v-btn
          color="error"
          text
          @click="
            () => {
              dialog = false;
              action();
            }
          "
        >
          Confirm
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
import Vue from "vue";

export default Vue.extend({
  props: {
    title: {
      type: String,
      default: "Suppression"
    },
    btnText: {
      type: String,
      default: "Supprimer"
    },
    btnSmall: {
      type: Boolean,
      default: false
    },
    icon: {
      type: Boolean,
      default: true
    },
    action: {
      type: Function,
      default: () => {}
    }
  },
  data() {
    return {
      dialog: false
    };
  },
  methods: {}
});
</script>
