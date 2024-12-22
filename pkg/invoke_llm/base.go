package invoke_llm

import "regexp"


var gitLinkPatterns = []*regexp.Regexp{
	regexp.MustCompile(`https?://github\.com/([A-Za-z0-9]+)/([A-Za-z0-9]+)`),
	regexp.MustCompile(`https?://gitlab\.com/[^,)#\s'"]+`),
	regexp.MustCompile(`https?://bitbucket\.org/[^,)#\s'"]+`),
	regexp.MustCompile(`https?://gitlab\.org/[^,)#\s'"]+`),
	regexp.MustCompile(`https?://gitee\.com/[^,)#\s'"]+`),
}

type GitLink struct {
	URL     string
	Pattern *regexp.Regexp
}

func CheckIfGitLink(url string) *GitLink {
	for _, pattern := range gitLinkPatterns {
		if pattern.MatchString(url) {
			return &GitLink{URL: url, Pattern: pattern}
		}
	}
	return nil
}

var PROMPT = map[string]string{
	"home2git_link": "Given the list of repository links [%s] and the homepage URL '%s', select the most likely repository link that matches the homepage. If a direct match is found, return the URL in the format 'URL is: [matched_url]'. If no direct match is found, check other repositories on GitHub, GitLab, or Gitee. If a related repository exists, respond with 'URL is: [url]'. If no relevant repository can be identified, respond with 'does not exist'.",
	"home2git_nolink": "Check if there is a git repository for %s hosted on platforms like GitHub, GitLab, or Gitee. If it exists, respond in the format 'URL is: [url]'. If no repository exists, respond with 'does not exist'.",
	"industry_idx": `
Next I'll give you a url that represents a github repository, the repository's readme, description, and topics, and you'll need to have a rough idea of the nature of the repository based on what you know about its contents, uses, etc. 
Then, based on your knowledge and judgment of this repository, categorize it into one of the following industry classifications step by step. 

url: %s
readme: %s
description: %s
topics: %s

Having given you the information about the warehouse above, you next have to choose the most relevant category from the following classifications
Classification and Introduction
Energy: The energy sector involves the production, distribution, and consumption of energy in various forms such as electricity, oil, gas, and renewable sources. It is crucial for powering industries, homes, and transportation, and plays a key role in economic development and environmental sustainability.
Transportation: The transportation sector includes the systems and means of moving people and goods by land, sea, and air. It is vital for economic growth, enabling trade, travel, and communication, and supports the functioning of modern society by connecting different regions and facilitating mobility.
Water Resources Management: Water resources management involves the planning, development, and conservation of water for various uses such as agriculture, industry, drinking water supply, and environmental preservation. It encompasses the management of rivers, lakes, reservoirs, and groundwater to ensure sustainable use and availability of water resources for present and future generations.
Finance: Finance refers to the management of money and investments, including banking, insurance, asset management, and financial markets. It plays a critical role in allocating resources, facilitating economic activities, and managing risks. The finance sector encompasses institutions, regulations, and instruments that enable individuals, businesses, and governments to handle financial transactions and achieve financial goals.
E-Government (Electronic Government): E-Government refers to the use of digital technologies, such as the internet and information and communication technology (ICT), by government agencies to enhance the delivery of public services, improve efficiency, and promote transparency and citizen participation. It involves the electronic exchange of information, communication of policies and services, and online interaction between government and citizens or businesses. E-Government initiatives aim to modernize governance processes, streamline administrative procedures, and enhance accessibility and responsiveness in public administration.
Defense Science and Technology Industry: The defense science and technology industry involves the research, development, production, and application of technologies and systems used for national defense and security purposes. It includes sectors such as aerospace, military electronics, weaponry, cybersecurity, and advanced materials. This industry plays a critical role in safeguarding national sovereignty, enhancing military capabilities, and maintaining strategic readiness in the face of evolving threats and geopolitical challenges.
Public Services: Public services refer to essential services provided by government or non-profit organizations to the general public. These services include healthcare, education, sanitation, transportation, and social welfare. They are crucial for ensuring the well-being and development of society, aiming to meet basic needs and promote equitable access to essential resources and opportunities.
Public Communication and Information Services: Public Communication and Information Services involve the infrastructure and services that facilitate the transmission of information across various platforms, including telecommunications, internet, and broadcasting. This sector is essential for connecting people and supporting economic and social activities by ensuring efficient and reliable communication.
General: If you think it meets the needs of various industries, then you need to choose General. 
Others: If you don't think it meets any of the above industry needs, then you need to choose Others.
[IMPORTANT] You were given 10 categories above and your answer must have one and only one of the above 10 categories. You only need to give me the result, and you must not give me any other explanation or other classification.
`,
}