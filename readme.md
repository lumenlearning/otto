# otto

Otto is named after Autolycus in the Shakespearean classic, Winter's Tale.
Autolycus is a thief, or as he puts it, "a snapper up of unconsidered trifles."

Otto was born out of the need to identify which rubric items on a particular
assignment also happened to be Learning Outcome objects.  At the time of 
writing, Canvas provided no reasonable means through the API to programmatically
discern whether a rubric item was in fact linked to a Learning Outcome.

Otto, therefore, does the following:

- Log in to an instance of Canvas using the supplied credentials
- Given a list of courses and assignments, retrieve the HTML for the rubric
  for each assignment
- Parse the HTML, looking for: <span class="learning_outcome_id"
- For each outcome, emit a CSV record (course_id, assignment_id, outcome_id).

The result is a CSV file showing all of the outcomes that were linked to the
various assignments in a course.  This information was necessary in order to
run the ImprovOER analysis.

# Usage

command line arguments for otto:

- -idfile="":
    The path to the file that contains the iDs of the courses and assignments
    we want to find outcomes for. The file is comma-separated (CSV) format
    (course_id,assignment_id).
- -login1="https://lumen.instructure.com/login":
    The first URL to load that will give us the needed authenticy_token
    information for the actual login.
- -login2="https://lumen.instructure.com/login?nonldap=true":
    The second URL to POST to in order to complete the login process.
- -username="somefakeuser@email":
    The username of the Canvas user account we'll be using for scraping.
- -password="somesecretpass!":
    The password for the user account.
- -useragent="Mozilla/5.0 (Macintosh; Intel Mac OS X 10_6_8) AppleWebKit/537.13+ (KHTML, like Gecko) Version/5.1.7 Safari/534.57.2":
    The User-Agent string to send in all HTTP requests.


# Example

otto -idfile=course_assignments.csv -username=somefakeuser@email -password="somesecretpass!"