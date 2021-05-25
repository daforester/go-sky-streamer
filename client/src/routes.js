const Main = () => import('views/Main.vue');

export const routes = [
  { path: '*', component: Main },
  { path: '/', component: Main },
];

export default routes;
