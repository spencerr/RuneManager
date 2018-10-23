<template lang="html">
    <div id="wrapper">
        <top-bar></top-bar>
        <left-side-bar></left-side-bar>
        <div class = "content-page">
            <div class = "content">
                <div class = "content-fluid">
                    Dashboard
                </div>
            </div>
        </div>
        <right-side-bar></right-side-bar>
    </div>
</template>


<script>
    import LeftSideBar from './LeftSideBar.vue'
    import RightSideBar from './RightSideBar.vue'
    import TopBar from './TopBar.vue'

    export default {
        components: {
            LeftSideBar, RightSideBar, TopBar
        }
    }
</script>