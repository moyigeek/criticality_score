
calc_score.rec: update_dependents.rec update_git_metrics.rec update_depsdev.rec
	# Calculate the score
	echo "* Calculating the score..."
	touch $@

update_dependents.rec: package_updated.src union_gitlink.rec
	# Update the dependents
	echo "* Updating the dependents..."
	touch $@

update_git_metrics.rec: git_updated.src union_gitlink.rec
	# Update the Git metrics
	echo "* Updating the Git metrics..."
	touch $@

update_depsdev.rec: depsdev_updated.src union_gitlink.rec
	# Update the dependencies development
	echo "* Updating the dependencies development..."
	touch $@


union_gitlink.rec: gitlink_updated.src
	# Union the Git links
	echo "* Unioning the Git links..."
	touch $@

