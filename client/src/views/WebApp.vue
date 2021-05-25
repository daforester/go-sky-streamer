<template>
  <div id="WebAppBody" class="WebAppBody">
    <router-view class="container-fluid" />
  </div>
</template>
<script>
import 'libs/fontawesome/minimal';

export default {
  name: 'WebApp',
};

</script>
<style lang="scss">
body {
  overflow-x: hidden;
  margin: 0;
  padding: 0;
}
#WebAppBody {
  @import '../libs/scss/variables.scss';
  @import '../libs/scss/controls.scss';

  background: white;
  color: #333;
  font-family: Segoe UI,
               Tahoma,
               Microsoft SansSerif,
               MS SansSerif,
               DejaVu Sans,
               Noto Sans,
               San Francisco,
               Roboto,
               Droid Sans,
               Helvetica,
               Arial,
               sans-serif;
}
</style>
