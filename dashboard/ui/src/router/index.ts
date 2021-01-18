import Vue from "vue";
import VueRouter, { RouteConfig } from "vue-router";
import Home from "../views/Home.vue";
import Medias from "../views/Medias.vue";
import NotFound from "../views/NotFound.vue";

Vue.use(VueRouter);

const routes: Array<RouteConfig> = [
  {
    path: "/",
    name: "Home",
    component: Home
  },
  {
    path: "/medias",
    name: "MediaAll",
    component: Medias
  },
  {
    path: "/medias/create",
    name: "MediaCreate",
    component: NotFound
  },
  {
    path: "/medias/:id",
    name: "MediaOne",
    component: NotFound
  },
  {
    path: "/metrics",
    name: "Metrics",
    component: NotFound
  },
  {
    path: "/settings",
    name: "Settings",
    component: NotFound
  },
  {
    path: "*",
    name: "NotFound",
    component: NotFound
  }
];

const router = new VueRouter({
  mode: "hash",
  base: process.env.BASE_URL,
  routes
});

export default router;
