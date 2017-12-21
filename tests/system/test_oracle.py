import os
import metricbeat
from nose.plugins.attrib import attr


class Test(metricbeat.BaseTest):

    def common_checks(self, output):
        # Ensure no errors or warnings exist in the log.
        log = self.get_log()
        self.assertNotRegexpMatches(log, "ERR|WARN")

        for evt in output:
            top_level_fields = metricbeat.COMMON_FIELDS + ["oracle"]
            self.assertItemsEqual(self.de_dot(top_level_fields), evt.keys())

            self.assert_fields_are_documented(evt)

    def get_hosts(self):
        return [os.getenv("ORACLE_DSN")]
