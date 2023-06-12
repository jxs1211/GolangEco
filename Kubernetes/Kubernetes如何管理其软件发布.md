# 版本更新不断，Kubernetes是如何管理其软件发布的？

CNCF [CNCF](javascript:void(0);) *2023-06-16 11:27* *Posted on 香港*

*作者：**Leonard Pahlke**[1]，CNCF 大使兼环境可永续发展 TAG 主席*

在本文中，我们将探讨开源项目 Kubernetes 如何管理其软件发布。通过探索已建立和演变的社区结构，本文将深入探讨支撑 Kubernetes 发布管理的协作动态，将重点从 Kubernetes 本身转移。

## 开源是令人着迷的 ✨

![Image](https://mmbiz.qpic.cn/mmbiz_png/GpkQxibjhkJyG1Eh4EIuDF6rjWZBYicuvNNCFxwpKTxSibtxwXS14lIMJv7icR8ibltAaRnxibd7Tibeqgic3G08m928Ag/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

开源无疑是令人着迷的。想想像 Kubernetes 这样的大型项目，全世界数千人汇聚在一起，共同构建正在改变整个行业的软件。贡献者们把自己的空闲时间和热情投入到这个项目中，同时还要应对自己的职业和个人责任。通过像 Kubernetes 这样的倡议，云原生领域已经出现，彻底改变了应用程序的开发和部署方式。

开源是创新的强大驱动力。通过拥抱开源，公司可以利用一个集体知识和专业知识的池，从协作努力中获得不断推出的进步和新的观点。忽视开源意味着可能错过创新和人才，这可能会使公司落后于其竞争对手。因此，积极参与开源项目不仅是保持与最新进展联系的一种方式，而且也是积极参与塑造技术未来的手段。

开源通过让不同的声音和思想在民主进程中汇聚，挑战了传统的等级制度和结构。软件开发中的民主化和包容原则可能促使公司采用类似的方法来开展其内部开发项目。在组织内建立协作、透明度和所有权文化，可以促进组织内更大的创造力和创新。

开源的吸引力不仅在于推动创新的能力，还在于它为个人和专业成长提供的机会。参与开源项目让人们能够从顶级专业人士中学习，并与一个充满激情、致力于推动技术边界的全球社区联系在一起。

## 开源软件项目中的挑战

参与开源软件项目很有趣，但也很重要的是要知道它们有自己的问题。开源软件开发的混乱性是其中最引人入胜的方面之一，也是一些障碍的来源。

- **协调和协作**：开源项目通常涉及一个遍布全球的遥远的贡献者社区。协调努力、设定目标并管理许多不同人的贡献可能会很困难。沟通和协作需要采取积极措施，以确保每个人都在同一页面上，朝着共同的目标努力。
- **缺乏集中控制**：在开源项目中，通常没有中央机构或控制对开发过程施加支配。相反，决策是通过讨论和共识建立集体做出的。这种基于共识的方法可能很耗时，并需要做出妥协。
- **社区管理**：开源项目在很大程度上依赖于社区参与和参与度。管理多样化的社区、鼓励参与以及培养积极和包容的环境需要有效的社区管理策略。强大的社区生态系统对项目的长期可永续性至关重要。

虽然在开源领域还有其他挑战，例如许可和知识产权，但上述列表应该足以提供一个充分的基础来进行下一步工作。解决这些挑战对于成功管理 Kubernetes 项目的版本发布至关重要。

## Kubernetes 社区结构：SIG-Town 比喻

![Image](https://mmbiz.qpic.cn/mmbiz_png/GpkQxibjhkJyG1Eh4EIuDF6rjWZBYicuvNzlicMxm7IGFtgpTe7BT0b7ADaevWDEWHfbQsOwXPTjDOB7SW9qM9zog/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1)

受 Paris Pittman 在 Google Kubernetes Podcast 的**近期一集**[2]中的深刻比喻启发，Kubernetes 社区可以比作一个城市。这个比喻引起了共鸣，并导致了 SIG-Town 图表的创建，展示了社区内部的组织结构。

Kubernetes 社区规模庞大且多样化。**CNCF Devstats**[3]数据突出了其规模，过去 12 个月中记录了近 3,000 名个人贡献者和大约 60,000 个提交。管理如此庞大而活跃的社区，推动决策，促进包容性并确保有效的协作可能是一个重大挑战。了解 Kubernetes 社区的不同 SIG、角色和流程非常重要。SIG-Town 比喻展示了 Kubernetes 社区有多么庞大、复杂和耐心。

## 协作

为了使开源项目彼此协调，制定各种规则和准则是非常重要的。这些规则有助于每个人以良好的方式协同工作。

### 跨界协作的概念

改善协作对于像 Kubernetes 这样的开源项目至关重要。以下是一些帮助人们在社区中协同工作的概念：

1. **一切都是文件**：每个贡献都通过版本控制系统 Git 进行管理。这意味着更改、讨论和决策过程围绕项目仓库中的文件展开。承认和尊重项目贡献者的贡献至关重要。但值得注意的是，在某些讨论中，很难将所有内容都保存在 markdown 文件或类似格式的文件中。在这种情况下，选择像公共 Google 文档文件或 GitHub 问题等替代平台会更适合并有助于有效的协作。
2. **所有权和维护者角色**：如果没有维护者承担所有权，社区将无法就决策达成共识并确保项目成熟。所有权在文件 OWNERS 中定义，列出了审核人和评审人。审核人和评审人负责特定领域的所有权，并贡献他们的专业知识来指导和审查贡献。除了代码审核人和评审人角色之外，还有其他角色，例如 SIG Chairs、Release Team Roles 等。担任这些角色意味着对社区有一定的承诺。
3. **能人主义**：基于展示的技能、知识和贡献赢得认可和影响力。这种方法强调了个人贡献的重要性，而不是层级职位或关联。
4. **透明度和开放性**：包括公开的沟通、决策过程和社区驱动的倡议。鼓励其他人参与，并为他们提供表达意见和关注的平台至关重要。重要的是要避免私人直接消息（direct message，DM），并确保所有沟通都在公开场合进行。

开源项目中的协作是一个广泛的主题，还有许多需要探索的领域。然而，了解这些关键概念为理解 Kubernetes 发布流程以及社区内的协作方式提供了坚实的基础。

### Kubernetes 增强提案（KEP）——“功能设计合同”

在 Kubernetes 项目中，通过创建 Kubernetes 增强提案（Kubernetes Enhancement Proposa，KEP）来促进对贡献的管理。KEP 作为一个设计文档，概述了要实现的期望功能。它充当一种合同协议，允许对所提出的更改进行评估、讨论和评估。对于对 Kubernetes 所做的任何修改，KEP 过程都是强制性的。简单地打开一个拉取请求（pull request，PR）来引入对新协议的支持或类似更改是不够的。为确保正确的文档，下面提供的**模板**[4]说明了 KEP 中需要包含的必要结构和信息。

```
Summary
Motivation
– Goals
– Non-Goals
Proposal
– User Stories (Optional)
– Notes/Constraints/Caveats (Optional)
– Risks and Mitigations
Design Details
– Test Plan
    – Prerequisite testing updates
    – Unit tests
    – Integration tests
    – e2e tests
– Graduation Criteria
– Upgrade / Downgrade Strategy
– Version Skew Strategy
Production Readiness Review Questionnaire
– Feature Enablement and Rollback
– Rollout, Upgrade and Rollback Planning
– Monitoring Requirements
– Dependencies
– Scalability
– Troubleshooting
Implementation History
Drawbacks
Alternatives
Infrastructure Needed (Optional)
```

KEP 遵循 Google 设计文档的概念和原则。与 Google 设计文档一样，KEP 作为在 Kubernetes 项目中提出和记录功能增强的结构化框架。它们为提出功能的设计、理由和实现细节提供系统化的方法。

![Image]()

每个 Kubernetes 特别兴趣小组（SIG）负责管理和监督其相应的 Kubernetes 增强提案（KEP）。这些 KEP 可以处于 SIG 的视野内的不同开发阶段。它们可以正在讨论之中，正在考虑是否在发布中包含（opted-in），或者已经合并并分类为 alpha、beta 或 stable，表示不同成熟度和准备就绪的不同级别。KEP 的大小可能会有很大的差异，但一般情况下，更倾向于较小的 KEP，以实现原子性的更改。

## SIG-Release

正如前面所提到的，管理和监控所有 Kubernetes 增强提案（KEP）和相关讨论是一个极具挑战性的任务。SIG-Release 承担了这一责任，其**宪章**[5]中概述了这一点。SIG-Release 的主要目标，是确保提供高质量的 Kubernetes 发布，不断增强发布和开发流程，并促进与下游社区的协作，这些社区利用 Kubernetes 发布来构建自己的制品。

### 发布团队结构

在**sig-release**[6]中，有两个团队：**发布工程团队**[7]负责开发发布工具，由更稳定的贡献者组成；**发布团队**[8]每个周期都会更改。发布团队的主要作用是促进发布流程。

![Image]()

参考上面的图表，团队由领导（Lead）和影子（Shadow）组成，他们在整个发布过程中负责不同的兴趣领域。影子计划作为 Sig-release、Kubernetes 社区和更广泛的开源领域的入门途径。你可以通过协助收集增强功能、监控测试管道等任务来为发布周期做出贡献。你可以在**此处**[9]了解有关发布团队角色和团队的更多信息。

对开源项目的贡献不仅仅是编码。你可以通过文档、社区参与、运营、治理、基础设施工作和其他各种领域的工作进行有价值的贡献。发布团队是一个典型的例子。

### 发布周期概述

典型的 Kubernetes 发布周期约为三个月，并在此期间涵盖了几个关键的截止日期。你可以在**此处**[10]找到 v1.28 发布周期的发布时间表。在过去的发布中，里程碑和截止日期没有发生重大变化，时间表相当一致。

- **组建发布团队**：发布团队开始工作，个人可以通过参与调查每个周期申请**加入**[11]团队。
- **开始增强跟踪**：SIG 在发布周期中选择他们的 KEP。发布团队确保选择的 KEP**满足所有要求**[12]，以便包含在发布周期中。
- **增强冻结**：需要在增强冻结时间之前合并 KEP。未及时合并的所有 KEP 都需要提出例外。有了增强冻结，发布团队就知道将进行哪些更改。
- **功能博客冻结**：KEP 作者可以选择表示他们有兴趣撰写有关其相关更改的博客文章。
- **代码冻结**：所有与已批准的 KEP 相关的代码更改必须成功合并。在代码冻结后，拉取请求将不会合并，并且需要提出例外，由发布负责人逐个评估。
- **主题截止日期**：标志着团队收集将在发布周期中进行的重大更改的截止日期。这些主要主题将突出显示在媒体中。
- **最终发布版本**：发布最终版本，发布博客，通过通信渠道宣布版本，更新文档并解除仓库限制。

每个 Kubernetes 发布周期都需要大量的工作。它涉及众多个人，多个截止日期和必须完成的任务。发布负责人负责维护整个发布周期的综合**跟踪问题**[13]，目前包括 100 多个复选框，代表各种“正式”任务。此外，特别是在沟通和增强发布团队的表现方面，需要大量额外的工作，以从一个周期到另一个周期进行迭代。

在我担任 1.26 版本的发布团队负责人期间，我清楚地记得收到来自 Slack 机器人的通知，我在一周内发送了超过 600 条消息，这清楚地表明了涉及的重大工作量。团队不断努力提高整体体验并尽可能自动化任务，但这样的改进需要时间和跨多个周期的迭代努力。

### 将 KEP 纳入发布流程：它是如何运作的？

![Image]()

确保将 KEP（Kubernetes 增强提案）包含在 Kubernetes 发布中涉及到两个关键方面需要注意。首先，KEP 作者必须承担的任务，其次，发布团队负责的截止日期和任务。

首先，KEP 作者必须与相关 SIG（特别兴趣小组）讨论 KEP。SIG 负责人在将 KEP 纳入发布中起着关键作用。KEP 作者必须勤奋地完成准备 KEP 以进行包含的所有必要任务。这包括获得 SIG 的批准和合并，打开拉取请求（PR）以更新代码，确保代码审查和批准，并更新文档。此外，KEP 作者可以选择撰写博客文章以提供更多见解。发布团队确保 KEP 作者完成了这些步骤，并密切监视过程，以防止负责 KEP 的 SIG 负责人错过任何事情。

### 为什么某些功能不包含在发布中？

![Image]()

重要的是要记住，大多数维护人员在业余时间为项目做出贡献。这是志愿工作，他们可以自由地在他们选择的事物上花费时间。

如果你注意到缺少某项功能，则有机会自己做出贡献。与工作中的工程经理联系，并解释这个功能对公司的 IT 系统的重要性。投资于 Kubernetes 项目不仅有利于你的公司，还允许你与才华横溢的工程师合作，并在这个过程中可能获得宝贵的知识。此外，你的贡献可以惠及更广泛的云原生社区，因为他们可能会发现你的更改非常有用。就是这样！如果你想与我联系，最好的方式是通过 Slack 的**Kubernetes**[14]或**CNCF**[15]工作区。

### 参考资料

[1]Leonard Pahlke: *https://twitter.com/leonardpahlke?lang=en*

[2]近期一集: *https://kubernetespodcast.com/episode/200-k8s-community-checkup/*

[3]CNCF Devstats: *https://all.devstats.cncf.io/d/53/projects-health-table?orgId=1*

[4]模板: *https://github.com/kubernetes/enhancements/blob/master/keps/NNNN-kep-template/README.md*

[5]宪章: *https://github.com/kubernetes/community/blob/master/sig-release/charter.md*

[6]sig-release: *https://github.com/kubernetes/sig-release*

[7]发布工程团队: *https://github.com/kubernetes/sig-release/tree/master/release-engineering*

[8]发布团队: *https://github.com/kubernetes/sig-release/tree/master/release-team*

[9]此处: *https://github.com/kubernetes/sig-release/tree/master/release-team*

[10]此处: *https://github.com/kubernetes/sig-release/tree/master/releases/release-1.28#timeline*

[11]加入: *https://github.com/kubernetes/sig-release/blob/master/release-team/release-team-selection.md*

[12]满足所有要求: *https://github.com/kubernetes/sig-release/blob/master/releases/release_phases.md#enhancements-freeze*

[13]跟踪问题: *https://github.com/kubernetes/sig-release/issues/2223*

[14]Kubernetes: *https://slack.k8s.io/*

[15]CNCF: *https://cloud-native.slack.com/*